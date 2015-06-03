#include <memory.h>
#include <stdint.h>

#include <stdio.h>

// Accelerates RgbaFindPillars.
void GoRgbaFindPillars(void* rgbaBytes, void* pillarBytes, int width,
    int height, int pillarCount, uint8_t minR, uint8_t minG, uint8_t minB,
    uint8_t maxR, uint8_t maxG, uint8_t maxB) {

  uint32_t *rgbaPixel = (uint32_t*)rgbaBytes;
  int *pillars = (int*)pillarBytes;

  memset(pillars, 0, sizeof(int32_t) * 4 * pillarCount);
  int32_t* minPillar = pillars;
  int32_t minHeight = 0;

  for (int x = 0; x < width; ++x) {
    int pillarHeight = 0;
    uint32_t* rgbaPixel = (uint32_t*)rgbaBytes + x;
    for (int y = 0; y < height; ++y, rgbaPixel += width) {
      uint32_t rgba = *rgbaPixel;
      uint8_t r = rgba & 0xff;
      uint8_t g = (rgba >> 8) & 0xff;
      uint8_t b = (rgba >> 16) & 0xff;

      if (r >= minR && g >= minG && b >= minB &&
          r <= maxR && g <= maxG && b <= maxB) {
        pillarHeight += 1;
      } else {
        if (pillarHeight > minHeight) {
          minPillar[0] = pillarHeight;
          minPillar[1] = x;
          minPillar[2] = y - pillarHeight;
          minPillar[3] = y - 1;
          minHeight = pillarHeight;

          int32_t* pillar = pillars;
          for (int i = pillarCount; i > 0; --i, pillar += 4) {
            int32_t pillarHeight = pillar[0];
            if (pillarHeight < minHeight) {
              minHeight = pillarHeight;
              minPillar = pillar;
            }
          }
        }
        pillarHeight = 0;
      }
    }
  }
}

// Accelerates GoRgbaFindPuddle.
int GoRgbaFindPuddle(void* rgbaBytes, void* puddleBytes, int width,
    int height, int startY, int maxPuddleSize,uint8_t minR,
    uint8_t minG, uint8_t minB, uint8_t maxR, uint8_t maxG, uint8_t maxB) {
  uint32_t* rgbaPixels = (uint32_t*)rgbaBytes;

  uint32_t* rgbaPixel0 = rgbaPixels + width * startY;
  for (int y0 = startY; y0 < height; ++y0) {
    for (int x0 = 0; x0 < width; ++x0, ++rgbaPixel0) {
      uint32_t rgba = *rgbaPixel0;
      uint8_t r = rgba & 0xff;
      uint8_t g = (rgba >> 8) & 0xff;
      uint8_t b = (rgba >> 16) & 0xff;
      if (r < minR || g < minG || b < minB ||
          r > maxR || g > maxG || b > maxB) {
        continue;
      }

      uint8_t a = (rgba >> 24) & 0xff;
      if (a == 0) {
        continue;
      }

      int32_t* puddleIn = (int32_t*)puddleBytes;
      int32_t* puddleOut = puddleIn;
      puddleOut[0] = x0;
      puddleOut[1] = y0;
      puddleOut += 2;
      *rgbaPixel0 &= 0x00ffffff;
      int puddleSize = 1;  // Found a puddle.
      if (puddleSize == maxPuddleSize) {
        return puddleSize;
      }
      while (puddleOut != puddleIn) {
        // NOTE: We're overwriting the for loop variables. This is OK because
        //       we know for sure that we're returning before the loop
        //       continues.
        x0 = puddleIn[0];
        y0 = puddleIn[1];
        puddleIn += 2;
        for (int dx = -1; dx <= 1; ++dx) {
          for (int dy = -1; dy <= 1; ++dy) {
            int x = dx + x0;
            int y = dy + y0;
            uint32_t *rgbaPixel = rgbaPixels + y * width + x;
            uint32_t rgba = *rgbaPixel;
            uint8_t r = rgba & 0xff;
            uint8_t g = (rgba >> 8) & 0xff;
            uint8_t b = (rgba >> 16) & 0xff;
            if (r < minR || g < minG || b < minB ||
                r > maxR || g > maxG || b > maxB) {
              continue;
            }
            uint8_t a = (rgba >> 24) & 0xff;
            if (a == 0) {
              continue;
            }

            puddleOut[0] = x;
            puddleOut[1] = y;
            puddleOut += 2;
            *rgbaPixel = rgba & 0x00ffffff;
            ++puddleSize;
            if (puddleSize == maxPuddleSize) {
              return puddleSize;
            }
          }
        }
      }
      return puddleSize;
    }
  }
  return 0;  // Did not find a puddle.
}

// Accelerates RgbaResetPuddles.
void GoRgbaResetPuddles(void* rgbaBytes, int byteSize) {
  uint8_t* alphaPixel = (uint8_t*)rgbaBytes + 3;
  for (int i = byteSize >> 2; i > 0; --i, alphaPixel += 4) {
    *alphaPixel = 255;
  }
}
