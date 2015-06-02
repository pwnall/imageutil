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
