#include <stdint.h>

#include <stdio.h>

// Accelerates MaskRgba.
void GoMaskRgba(void* bytes, int byteCount, uint64_t mask) {
  uint64_t* words = (uint64_t*)bytes;
  for (int wordCount = (byteCount >> 3); wordCount > 0; --wordCount) {
    *words &= mask;
    ++words;
  }
}

// Accelerates RgbaToHsla.
void GoRgbaToHsla(void* rgbaBytes, void* hslaBytes, int byteCount) {
  uint8_t *rgbaPixel = (uint8_t*)rgbaBytes;
  uint8_t *hslaPixel = (uint8_t*)hslaBytes;
  for (int i = byteCount; i > 0; i -= 4, rgbaPixel += 4, hslaPixel += 4) {
    int r = rgbaPixel[0];
    int g = rgbaPixel[1];
    int b = rgbaPixel[2];
    int a = rgbaPixel[3];

    // RGB -> HSL formula lifted and adapted for 0-255 from:
    // http://www.niwa.nu/2013/05/math-behind-colorspace-conversions-rgb-hsl/
    int min = r;
    if (min > g) min = g;
    if (min > b) min = b;
    int max = r;
    if (max < g) max = g;
    if (max < b) max = b;

    int sum = min + max;
    int diff = max - min;
    int l = sum >> 1;  // L = (min + max) / 2
    int h, s;
    if (diff == 0) {
      s = 0;
      h = 0;
    } else {
      if (l >= 128) {
        s = 255 * diff / (510 - sum);
      } else {
        s = 255 * diff / sum;
      }

      if (max == r) {
        h = 42 * (g - b) / diff;
      } else if (max == g) {
        h = 84 + 42 * (b - r) / diff;
      } else {
        h = 168 + 42 * (r - g) / diff;
      }
      if (h < 0) h += 256;
    }

    hslaPixel[0] = (uint8_t)h;
    hslaPixel[1] = (uint8_t)s;
    hslaPixel[2] = (uint8_t)l;
    hslaPixel[3] = a;
  }
}

// Accelerates RgbaThreshold.
void GoRgbaThreshold(void* rgbaBytes, int byteCount, uint8_t minR,
    uint8_t minG, uint8_t minB, uint8_t maxR, uint8_t maxG, uint8_t maxB) {
  uint8_t *rgbaPixel = (uint8_t*)rgbaBytes;
  for (int i = byteCount; i > 0; i -= 4, rgbaPixel += 4) {
    uint8_t r = rgbaPixel[0];
    uint8_t g = rgbaPixel[1];
    uint8_t b = rgbaPixel[2];

    uint8_t a;
    if (r >= minR && g >= minG && b >= minB &&
        r <= maxR && g <= maxG && b <= maxB) {
      a = 255;
    } else {
      a = 0;
    }
    rgbaPixel[3] = a;
  }
}
