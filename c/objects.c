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
