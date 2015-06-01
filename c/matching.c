#include <memory.h>
#include <stdint.h>

// Accelerates RgbaCheckCrop.
int GoRgbaCheckCrop(void* haystackBytes, int hayWidth, void* needleBytes,
    int needleWidth, int needleHeight, int needleX, int needleY) {
  int haystackStride = 4 * hayWidth;
  uint32_t* haystackRow = (uint32_t*)haystackBytes +
      needleY * haystackStride + needleX;

  int needleStride = 4 * needleWidth;
  uint32_t* needleRow = (uint32_t*)needleBytes;

  for (int y = needleHeight; y > 0; --y) {
    if (!memcmp(haystackRow, needleRow, needleStride))
      return 0;
    haystackRow += haystackStride;
    needleRow += needleStride;
  }
  return 1;
}

// Accelerates RabinKarp.
void GoRabinKarp(void* worldBytes, void* scratch, uint32_t targetHash,
    int width, int height, int worldWidth, int worldHeight) {
  uint32_t* worldPixels = (uint32_t*)worldBytes;
  uint32_t* chash = (uint32_t*)scratch;  // column hashes
  const uint32_t kx = 1000000007;
  const uint32_t ky = 1000000007;
  const uint32_t m = 2000000011;

  uint32_t kx_w = 1;  // kx ^ w % m
  for (int x = 0; x < width; ++x) {
    kx_w = (uint32_t)(((uint64_t)kx_w * kx) % m);
  }
  uint32_t ky_h = 1;  // ky ^ h % m
  for (int y = 0; y < height; ++y) {
    ky_h = (uint32_t)(((uint64_t)ky_h * ky) % m);
  }

  memset(chash, 0, sizeof(uint32_t) * worldWidth);
  for (int y = 0; y < height - 1; ++y) {
    uint32_t* row = &worldPixels[y * worldWidth];
    uint32_t hash0y = 0;
    for (int x = 0; x < worldWidth; ++x) {
      chash[x] = (uint32_t)(((uint64_t)chash[x] * ky + row[x]) % m);
    }
  }

  // Special-case for first row.
  {
    int x = 0;
    uint32_t hash = 0;
    uint32_t* row = &worldPixels[(height - 1) * worldWidth];
    for (; x < width; ++x) {
      hash = (uint32_t)(((uint64_t)hash * kx + row[x]) % m);
    }
    if (hash == targetHash) {
      // TODO: check match at 0, 0
    }
    for (; x <  worldWidth; ++x) {
      hash = (uint32_t)(((uint64_t)hash * kx + row[x]) % m);
      hash = (m + hash -
          (uint32_t)(((uint64_t)kx_w * row[x - width]) % m)) % m;
      if (hash == targetHash) {
        // TODO: check match at x - width, 0
      }
    }
  }

  for (int y = height; y < worldHeight; ++y) {
    uint32_t* row = &worldPixels[y * worldWidth];
    uint32_t* oldRow = &worldPixels[(y - height) * worldWidth];
    uint32_t hash = 0;
    for (int x = 0; x < width; ++x) {

    }

    for (int x = width; x < worldWidth; ++x) {

    }
  }

  uint32_t whash;
}
