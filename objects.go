package imageutil

// #include "c/objects.c"
import "C"  // cgo

import (
  "unsafe"
)

// RgbaFindPillars returns the tallest vertical strips in an image.
func RgbaFindPillars(rgbaImage []byte, width int, height int,
    minRed int, maxRed int, minGreen int, maxGreen int, minBlue int,
    maxBlue int, pillars [][4]int32) {
  if cap(rgbaImage) < 4 * width * height {
    panic("RGBA image capacity inconsistent with width / height")
  }
  C.GoRgbaFindPillars(unsafe.Pointer(&rgbaImage[0]),
      unsafe.Pointer(&pillars[0][0]), C.int(width), C.int(height),
      C.int(len(pillars)), C.uint8_t(minRed), C.uint8_t(minGreen),
      C.uint8_t(minBlue), C.uint8_t(maxRed), C.uint8_t(maxGreen),
      C.uint8_t(maxBlue))
}

// RgbaFindPuddle locates contiguous areas in an image.
// The image's A channel is (ab)used to track the image's visited areas.
// It returns the size of the area that it found.
func RgbaFindPuddle(rgbaImage []byte, width int, height int,
    minRed int, maxRed int, minGreen int, maxGreen int, minBlue int,
    maxBlue int, startY int, puddlePixels [][2]int32) int {
  if cap(rgbaImage) < 4 * width * height {
    panic("RGBA image capacity inconsistent with width / height")
  }
  result := C.GoRgbaFindPuddle(unsafe.Pointer(&rgbaImage[0]),
      unsafe.Pointer(&puddlePixels[0][0]), C.int(width), C.int(height),
      C.int(startY), C.int(len(puddlePixels)), C.uint8_t(minRed),
      C.uint8_t(minGreen), C.uint8_t(minBlue), C.uint8_t(maxRed),
      C.uint8_t(maxGreen), C.uint8_t(maxBlue))
  return int(result)
}

// RgbaResetPuddles resets the Alpha channel of all pixles to 255.
// This is useful after running puddle searches over an image.
func RgbaResetPuddles(rgbaImage []byte) {
  C.GoRgbaResetPuddles(unsafe.Pointer(&rgbaImage[0]), C.int(len(rgbaImage)))
}
