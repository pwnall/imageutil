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
void GoRgbaToHsla(void* rgbaBytes, void* hslaBytes, int byteCount) {
  uint32_t *rgbaPixel = (uint32_t*)rgbaBytes;
  uint32_t *hslaPixel = (uint32_t*)hslaBytes;
  for (int i = (byteCount >> 2); i > 0; --i, ++rgbaPixel, ++hslaPixel) {
    unsigned rgba = *rgbaPixel;
    int r = rgba & 0xff;
    int g = (rgba >> 8) & 0xff;
    int b = (rgba >> 16) & 0xff;
    unsigned unshiftedA = rgba & 0xff000000;

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
        if (h < 0) h += 256;
      } else if (max == g) {
        h = 84 + 42 * (b - r) / diff;
      } else {
        h = 168 + 42 * (r - g) / diff;
      }
    }

    uint32_t hsla = (uint32_t)((unsigned)h | ((unsigned)s << 8) |
        ((unsigned)l << 16) | unshiftedA);
    *hslaPixel = hsla;
  }
}

// Accelerates RgbaThreshold.
void GoRgbaThreshold(void* rgbaBytes, int byteCount, uint8_t minR,
    uint8_t minG, uint8_t minB, uint8_t maxR, uint8_t maxG, uint8_t maxB) {
  uint32_t *rgbaPixel = (uint32_t*)rgbaBytes;
  for (int i = (byteCount >> 2); i > 0; --i, ++rgbaPixel) {
    unsigned rgba = *rgbaPixel & 0x00ffffff;
    uint8_t r = rgba & 0xff;
    uint8_t g = (rgba >> 8) & 0xff;
    uint8_t b = rgba >> 16;

    if (r >= minR && g >= minG && b >= minB &&
        r <= maxR && g <= maxG && b <= maxB) {
      rgba |= 0xff000000;
    }
    *rgbaPixel = rgba;
  }
}
