package imageutil

// #include "c/filters.c"
import "C"  // cgo

import (
  "unsafe"
)

// BuildMask computes a word mask from a 32-bit RGBA mask.
// The word mask is intended to be used with the MaskRgba function.
func BuildRgbaMask(rgba uint32) uint64 {
  // NOTE: Intels are little-endian, so we need to flip the bytes in the word.
  argb := uint64(((rgba & 0xff) << 24) | ((rgba & 0xff00) << 8) |
      ((rgba & 0xff0000) >> 8) | ((rgba & 0xff000000) >> 24))

  return argb | (argb << 32)
}

// MaskRgba applies a word mask to an RGBA image buffer.
// This uses fast 64-bit operations. In return for the speed, the caller must
// covert the RGBA mask into a word mask, with the help of BuildMask.
func MaskRgba(rawImage []byte, mask uint64) {
  C.GoMaskRgba(unsafe.Pointer(&rawImage[0]), C.int(len(rawImage)),
      C.uint64_t(mask))
}

// RgbaToHsla converts an RGBA image to a HSLA image.
// H, S, and L are in the range 0..255. A is unchanged.
// HSLA images can be used with most RGBA processing functions, like RgbaToPng.
func RgbaToHsla(rgbaImage []byte, hslaImage []byte) {
  if (cap(hslaImage) < len(rgbaImage)) {
    panic("HSLA buffer smaller than RGBA image size")
  }
  C.GoRgbaToHsla(unsafe.Pointer(&rgbaImage[0]), unsafe.Pointer(&hslaImage[0]),
      C.int(len(rgbaImage)));
}

// RgbPixelToHsl returns the HSL values for a RGB color with 8-bits / channel.
// It is mostly useful for easy conversion to our custom HSL scheme where H
// is scaled between 0 and 255.
func RgbPixelToHsl(red int, green int, blue int) (int, int, int) {
  argb := uint32(uint32(red) | uint32(green << 8) | uint32(blue << 16))
  var alsh uint32
  C.GoRgbaToHsla(unsafe.Pointer(&argb), unsafe.Pointer(&alsh), C.int(4))

  h := int(alsh & 0xff)
  s := int((alsh >> 8) & 0xff)
  l := int((alsh >> 16) & 0xff)
  return h, s, l
}

// RgbaThreshold sets the alpha channel in image to a threshold function.
// The function is 1 when the R, G, and B values are between given amounts, and
// 0 otherwise.
func RgbaThreshold(rgbaImage []byte, minRed int, maxRed int, minGreen int,
    maxGreen int, minBlue int, maxBlue int) {
  C.GoRgbaThreshold(unsafe.Pointer(&rgbaImage[0]), C.int(len(rgbaImage)),
      C.uint8_t(minRed), C.uint8_t(minGreen), C.uint8_t(minBlue),
      C.uint8_t(maxRed), C.uint8_t(maxGreen), C.uint8_t(maxBlue))
}
