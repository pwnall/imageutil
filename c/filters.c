#include <stdint.h>

// Accelerates MaskRgba.
void GoMaskRgba(void* bytes, int byteCount, uint64_t mask) {
  uint64_t* words = (uint64_t*)bytes;
  for (int wordCount = (byteCount >> 3); wordCount > 0; --wordCount) {
    *words &= mask;
    ++words;
  }
}

// Accelerates RgbaToHsla.
void GoRgbaToHsla(void* rgbaBytes, void* hslaBytes, int pixelCount) {
  uint8_t *rgbaPixel = (uint8_t*)rgbaBytes;
  uint8_t *hslaPixel = (uint8_t*)hslaBytes;
  for (int i = pixelCount; i > 0; --i, rgbaPixel += 4, hslaPixel += 4) {
    int r = rgbaPixel[0];
    int g = rgbaPixel[1];
    int b = rgbaPixel[2];
    int a = rgbaPixel[3];

    int min = r;
    if (min > g) min = g;
    if (min > b) min = b;
    int max = r;
    if (max > g) max = g;
    if (max > b) max = b;

    // L = (min + max) / 2
    int l = (min + max) >> 1;
    int s;
    if (min == max) {
      s = 0;
    } else {
      if (l >= 128) {
        s = 255 * (max - min) / (510 - max - min);
      } else {
        s = 255 * (max - min) / (max + min);
      }
    }

    int h;
    if (max == r) {
      h = 42 * (g - b) / (max - min);
    } else if (max == g) {
      h = 84 + 42 * (b - r) / (max - min);
    } else {
      h = 168 + 42 * (r - g) / (max - min);
    }
    if (h < 0) h += 256;

    hslaPixel[0] = h;
    hslaPixel[1] = s;
    hslaPixel[2] = l;
    hslaPixel[3] = a;
  }
}
