package imageutil

// #include "c/matching.c"
import "C"  // cgo

import (
  "unsafe"
)

// RgbaCheckCrop returns true if an image is a cropped version of another one.
// This is image pattern-matching, but only aligns the pattern with the image
// in one predetermined position.
func RgbaCheckCrop(haystack []byte, hayWidth int, hayHeight int,
    needle []byte, needleWidth int, needleHeight int, needleX int,
    needleY int) bool {

  // NOTE: These checks are mainly here to prevent segmentation faults in the
  //       C code. Therefore, panicing is appropriate.
  if len(haystack) != hayWidth * hayHeight * 4 {
    panic("Haystack width and height do not match buffer size")
  }
  if len(needle) != needleWidth * needleHeight * 4 {
    panic("Needle width and height do not match buffer size")
  }

  // NOTE: These checks are also intended to prevent segmentation faults, but
  //       we don't have to panic here.
  if needleX < 0 || needleX + needleWidth > hayWidth {
    return false
  }
  if needleY < 0 || needleY + needleHeight > hayHeight {
    return false
  }

  // NOTE: The haystack's height is irrelevant to the actual matching logic,
  //       so it is omitted.
  cresult := C.GoRgbaCheckCrop(unsafe.Pointer(&haystack[0]), C.int(hayWidth),
      unsafe.Pointer(&needle[0]), C.int(needleWidth), C.int(needleHeight),
      C.int(needleX), C.int(needleY))
  return cresult != 0
}
