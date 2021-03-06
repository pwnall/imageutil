#include <memory.h>
#include <stdint.h>

// Accelerates RgbaCheckCrop.
int GoRgbaCheckCrop(void* haystackBytes, void* needleBytes, int hayWidth,
    int needleWidth, int needleHeight, int needleLeft, int needleTop) {
  uint32_t* haystackRow = (uint32_t*)haystackBytes + needleTop * hayWidth +
      needleLeft;
  uint32_t* needleRow = (uint32_t*)needleBytes;
  uint32_t needleStride = needleWidth * 4;
  for (int y = needleHeight; y > 0; --y) {
    if (memcmp(haystackRow, needleRow, needleStride))
      return 0;
    haystackRow += hayWidth;
    needleRow += needleWidth;
  }
  return 1;
}

// Accelerates RgbaCheckMaskedCrop.
int GoRgbaCheckMaskedCrop(void* haystackBytes, void* needleBytes, int hayWidth,
    int needleWidth, int needleHeight, int needleLeft, int needleTop,
    uint32_t argbMask) {
  uint32_t* haystackPtr = (uint32_t*)haystackBytes + needleTop * hayWidth +
      needleLeft;
  uint32_t* needlePtr = (uint32_t*)needleBytes;
  int rowJump = hayWidth - needleWidth;
  for (int y = needleHeight; y > 0; --y) {
    for (int x = needleWidth; x > 0; --x, ++needlePtr, ++haystackPtr) {
      if ((*haystackPtr & argbMask) != *needlePtr)
        return 0;
    }
    haystackPtr += rowJump;
  }
  return 1;
}

// Accelerates RgbaDiffMaskedCrop.
int64_t GoRgbaDiffMaskedCrop(void* haystackBytes, void* needleBytes,
    int hayWidth, int needleWidth, int needleHeight, int needleLeft,
    int needleTop, uint32_t argbMask) {
  uint32_t* haystackPtr = (uint32_t*)haystackBytes + needleTop * hayWidth +
      needleLeft;
  uint32_t* needlePtr = (uint32_t*)needleBytes;
  int rowJump = hayWidth - needleWidth;
  int64_t diff = 0;
  for (int y = needleHeight; y > 0; --y) {
    for (int x = needleWidth; x > 0; --x, ++needlePtr, ++haystackPtr) {
      uint32_t hrgba = (*haystackPtr & argbMask);
      uint32_t nrgba = *needlePtr;
      int hchannel, nchannel;

      // Red.
      hchannel = hrgba & 0xff; hrgba >>= 8;
      nchannel = nrgba & 0xff; nrgba >>= 8;
      diff += (hchannel >= nchannel) ?
          hchannel - nchannel : nchannel - hchannel;

      // Green.
      hchannel = hrgba & 0xff; hrgba >>= 8;
      nchannel = nrgba & 0xff; nrgba >>= 8;
      diff += (hchannel >= nchannel) ?
          hchannel - nchannel : nchannel - hchannel;

      // Blue.
      hchannel = hrgba & 0xff; hrgba >>= 8;
      nchannel = nrgba & 0xff; nrgba >>= 8;
      diff += (hchannel >= nchannel) ?
          hchannel - nchannel : nchannel - hchannel;

      // Alpha.
      hchannel = hrgba;
      nchannel = nrgba;
      diff += (hchannel >= nchannel) ?
          hchannel - nchannel : nchannel - hchannel;
    }
    haystackPtr += rowJump;
  }
  return diff;
}

// Accelerates RgbaDiffThresholdCrop.
int GoRgbaDiffThresholdCrop(void* haystackBytes, void* needleBytes,
    int hayWidth, int needleWidth, int needleHeight, int needleLeft,
    int needleTop,  uint8_t minR, uint8_t minG, uint8_t minB, uint8_t maxR,
    uint8_t maxG, uint8_t maxB) {
  uint32_t* haystackPtr = (uint32_t*)haystackBytes + needleTop * hayWidth +
      needleLeft;
  uint32_t* needlePtr = (uint32_t*)needleBytes;
  int rowJump = hayWidth - needleWidth;
  int diff = 0;
  for (int y = needleHeight; y > 0; --y) {
    for (int x = needleWidth; x > 0; --x, ++needlePtr, ++haystackPtr) {
      uint32_t hrgba = *haystackPtr;
      uint8_t hr = hrgba & 0xff; hrgba >>= 8;
      uint8_t hg = hrgba & 0xff; hrgba >>= 8;
      uint8_t hb = hrgba & 0xff;
      uint8_t ha = (hr >= minR && hg >= minG && hb >= minB &&
        hr <= maxR && hg <= maxG && hb <= maxB) ? 0xff : 0;

      uint32_t nrgba = *needlePtr;
      uint8_t nr = nrgba & 0xff; nrgba >>= 8;
      uint8_t ng = nrgba & 0xff; nrgba >>= 8;
      uint8_t nb = nrgba & 0xff; nrgba >>= 8;
      uint8_t na = nrgba;

      if (ha != na)
        diff += 1;
    }
    haystackPtr += rowJump;
  }
  return diff;
}



// (a * b) % m
static inline uint32_t mulMod(uint32_t a, uint32_t b, uint32_t m) {
  return (uint32_t)(((uint64_t)a * b) % m);
}
// (a - b) % m
static inline uint32_t modSub(uint32_t a, uint32_t b, uint32_t m) {
  return (uint32_t)(((uint64_t)m + a - b) % m);
}
// (a * b + c) % m
static inline uint32_t mulModAdd(uint32_t a, uint32_t b, uint32_t c,
    uint32_t m) {
  // NOTE: This doesn't overflow.
  //       (2^32 - 1) * (2^32 - 1) = 2^64 - 2 * 2^32 + 1,
  //       So we have room for 2 * 2^32 - 1 before we wrap around.
  return (uint32_t)(((uint64_t)a * b + c) % m);
}

// Multiplicative constant across column hashes.
static const uint32_t kx = 1000000007;
// Multiplicative constant across a column.
static const uint32_t ky = 1000000007;
// Hash modulo.
static const uint32_t m = 2000000011;

// Accelerates RabinKarpHash.
// This doesn't really need accelerating, but it's easier to just reuse the
// code in GoRabinKarp below and keep it in sync than to rewrite the whole
// thing in Go.
uint32_t GoHashForRgbaFindCrop(void *needleBytes, int needleWidth,
    int needleHeight) {
  uint32_t hash = 0;
  for (int x = 0; x < needleWidth; ++x) {
    uint32_t* column = (uint32_t*)needleBytes + x;
    uint32_t chash = 0;
    for (int y = 0; y < needleHeight; ++y) {
      chash = mulModAdd(chash, ky, *column, m);
      column += needleWidth;
    }
    hash = mulModAdd(hash, kx, chash, m);
  }
  return hash;
}

// Accelerates RgbaFindCrop.
// The scratch space must point to a buffer of worldWidth uint32_t elements.
int GoRgbaFindCrop(void* haystackBytes, void *needleBytes, int hayWidth,
    int hayHeight, int needleWidth, int needleHeight, uint32_t needleHash,
    void* scratch, int* matchLeft, int* matchTop) {
  uint32_t* hayPixels = (uint32_t*)haystackBytes;
  uint32_t* chash = (uint32_t*)scratch;  // column hashes

  uint32_t kx_w = 1;  // kx ^ w % m
  for (int x = 0; x < needleWidth; ++x) {
    kx_w = (uint32_t)(((uint64_t)kx_w * kx) % m);
  }
  uint32_t ky_h = 1;  // ky ^ h % m
  for (int y = 0; y < needleHeight; ++y) {
    ky_h = (uint32_t)(((uint64_t)ky_h * ky) % m);
  }

  memset(chash, 0, sizeof(uint32_t) * hayWidth);
  for (int y = 0; y < needleHeight; ++y) {
    uint32_t* row = &hayPixels[y * hayWidth];
    for (int x = 0; x < hayWidth; ++x) {
      chash[x] = mulModAdd(chash[x], ky, row[x], m);
    }
  }

  int matchCount = 0;

  // Special-case for first row.
  uint32_t hash = 0;
  {
    int x = 0;
    for (; x < needleWidth; ++x) {
      hash = mulModAdd(hash, kx, chash[x], m);
    }
    if (hash == needleHash) {
      int needleLeft = 0;
      int needleTop = 0;
      if (GoRgbaCheckCrop(haystackBytes, needleBytes, hayWidth, needleWidth,
            needleHeight, needleLeft, needleTop)) {
        matchCount += 1;
        *matchLeft = needleLeft;
        *matchTop = needleTop;
      }
    }

    for (; x < hayWidth; ++x) {
      hash = mulModAdd(hash, kx, chash[x], m);
      hash = modSub(hash, mulMod(chash[x - needleWidth], kx_w, m), m);

      if (hash == needleHash) {
        int needleLeft = x - needleWidth + 1;
        int needleTop = 0;
        if (GoRgbaCheckCrop(haystackBytes, needleBytes, hayWidth, needleWidth,
              needleHeight, needleLeft, needleTop)) {
          matchCount += 1;
          *matchLeft = needleLeft;
          *matchTop = needleTop;
        }
      }
    }
  }

  for (int y = needleHeight; y < hayHeight; ++y) {
    uint32_t* row = &hayPixels[y * hayWidth];
    uint32_t* oldRow = &hayPixels[(y - needleHeight) * hayWidth];
    hash = 0;
    for (int x = 0; x < needleWidth; ++x) {
      chash[x] = mulModAdd(chash[x], ky, row[x], m);
      chash[x] = modSub(chash[x], mulMod((uint64_t)oldRow[x], ky_h, m), m);

      hash = mulModAdd(hash, kx, chash[x], m);
    }

    if (hash == needleHash) {
      int needleLeft = 0;
      int needleTop = y - needleHeight + 1;
      if (GoRgbaCheckCrop(haystackBytes, needleBytes, hayWidth, needleWidth,
            needleHeight, needleLeft, needleTop)) {
        matchCount += 1;
        *matchLeft = needleLeft;
        *matchTop = needleTop;
      }
    }

    for (int x = needleWidth; x < hayWidth; ++x) {
      chash[x] = mulModAdd(chash[x], ky, row[x], m);
      chash[x] = modSub(chash[x], mulMod((uint64_t)oldRow[x], ky_h, m), m);

      hash = mulModAdd(hash, kx, chash[x], m);
      hash = modSub(hash, mulMod(chash[x - needleWidth], kx_w, m), m);
      if (hash == needleHash) {
        int needleLeft = x - needleWidth + 1;
        int needleTop = y - needleHeight + 1;
        if (GoRgbaCheckCrop(haystackBytes, needleBytes, hayWidth, needleWidth,
              needleHeight, needleLeft, needleTop)) {
          matchCount += 1;
          *matchLeft = needleLeft;
          *matchTop = needleTop;
        }
      }
    }
  }

  return matchCount;
}

// Accelerates RgbaFindMaskedCrop.
// The scratch space must point to a buffer of worldWidth uint32_t elements.
int GoRgbaFindMaskedCrop(void* haystackBytes, void *needleBytes, int hayWidth,
    int hayHeight, int needleWidth, int needleHeight, uint32_t argbMask,
    uint32_t needleHash, void* scratch, int* matchLeft, int* matchTop) {
  uint32_t* hayPixels = (uint32_t*)haystackBytes;
  uint32_t* chash = (uint32_t*)scratch;  // column hashes

  uint32_t kx_w = 1;  // kx ^ w % m
  for (int x = 0; x < needleWidth; ++x) {
    kx_w = (uint32_t)(((uint64_t)kx_w * kx) % m);
  }
  uint32_t ky_h = 1;  // ky ^ h % m
  for (int y = 0; y < needleHeight; ++y) {
    ky_h = (uint32_t)(((uint64_t)ky_h * ky) % m);
  }

  memset(chash, 0, sizeof(uint32_t) * hayWidth);
  for (int y = 0; y < needleHeight; ++y) {
    uint32_t* row = &hayPixels[y * hayWidth];
    for (int x = 0; x < hayWidth; ++x) {
      chash[x] = mulModAdd(chash[x], ky, row[x] & argbMask, m);
    }
  }

  int matchCount = 0;

  // Special-case for first row.
  uint32_t hash = 0;
  {
    int x = 0;
    for (; x < needleWidth; ++x) {
      hash = mulModAdd(hash, kx, chash[x], m);
    }
    if (hash == needleHash) {
      int needleLeft = 0;
      int needleTop = 0;
      if (GoRgbaCheckMaskedCrop(haystackBytes, needleBytes, hayWidth,
            needleWidth, needleHeight, needleLeft, needleTop, argbMask)) {
        matchCount += 1;
        *matchLeft = needleLeft;
        *matchTop = needleTop;
      }
    }

    for (; x < hayWidth; ++x) {
      hash = mulModAdd(hash, kx, chash[x], m);
      hash = modSub(hash, mulMod(chash[x - needleWidth], kx_w, m), m);

      if (hash == needleHash) {
        int needleLeft = x - needleWidth + 1;
        int needleTop = 0;
        if (GoRgbaCheckMaskedCrop(haystackBytes, needleBytes, hayWidth,
              needleWidth, needleHeight, needleLeft, needleTop, argbMask)) {
          matchCount += 1;
          *matchLeft = needleLeft;
          *matchTop = needleTop;
        }
      }
    }
  }

  for (int y = needleHeight; y < hayHeight; ++y) {
    uint32_t* row = &hayPixels[y * hayWidth];
    uint32_t* oldRow = &hayPixels[(y - needleHeight) * hayWidth];
    hash = 0;
    for (int x = 0; x < needleWidth; ++x) {
      chash[x] = mulModAdd(chash[x], ky, row[x] & argbMask, m);
      chash[x] = modSub(chash[x],
          mulMod((uint64_t)oldRow[x] & argbMask, ky_h, m), m);

      hash = mulModAdd(hash, kx, chash[x], m);
    }

    if (hash == needleHash) {
      int needleLeft = 0;
      int needleTop = y - needleHeight + 1;
      if (GoRgbaCheckMaskedCrop(haystackBytes, needleBytes, hayWidth,
            needleWidth, needleHeight, needleLeft, needleTop, argbMask)) {
        matchCount += 1;
        *matchLeft = needleLeft;
        *matchTop = needleTop;
      }
    }

    for (int x = needleWidth; x < hayWidth; ++x) {
      chash[x] = mulModAdd(chash[x], ky, row[x] & argbMask, m);
      chash[x] = modSub(chash[x],
          mulMod((uint64_t)oldRow[x] & argbMask, ky_h, m), m);

      hash = mulModAdd(hash, kx, chash[x], m);
      hash = modSub(hash, mulMod(chash[x - needleWidth], kx_w, m), m);
      if (hash == needleHash) {
        int needleLeft = x - needleWidth + 1;
        int needleTop = y - needleHeight + 1;
        if (GoRgbaCheckMaskedCrop(haystackBytes, needleBytes, hayWidth,
              needleWidth, needleHeight, needleLeft, needleTop, argbMask)) {
          matchCount += 1;
          *matchLeft = needleLeft;
          *matchTop = needleTop;
        }
      }
    }
  }

  return matchCount;
}
