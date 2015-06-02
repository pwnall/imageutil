package imageutil

// #include "c/objects.c"
import "C"  // cgo

import (
  "unsafe"
)

// RgbaFindPillars returns the tallest vertical strips in an image.
func RgbaFindPillars(rgbaImage []byte, width int, height int,
    minRed int, maxRed int, minGreen int, maxGreen int, minBlue int,
    maxBlue int, pillarCount int, pillars [][4]int32) {
  if cap(pillars) < pillarCount {
    panic("Insufficient pillar slice capacity")
  }
  C.GoRgbaFindPillars(unsafe.Pointer(&rgbaImage[0]),
      unsafe.Pointer(&pillars[0][0]), C.int(width), C.int(height),
      C.int(pillarCount), C.uint8_t(minRed), C.uint8_t(minGreen),
      C.uint8_t(minBlue), C.uint8_t(maxRed), C.uint8_t(maxGreen),
      C.uint8_t(maxBlue))
}
