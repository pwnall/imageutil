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
  if len(haystack) < hayWidth * hayHeight * 4 {
    panic("Haystack width and height do not match buffer size")
  }
  if len(needle) < needleWidth * needleHeight * 4 {
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
  cresult := C.GoRgbaCheckCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(needleWidth),
      C.int(needleHeight), C.int(needleX), C.int(needleY))
  return cresult != 0
}

// HashForRgbaFind computes the needle hash needed by RgbaFind.
// It returns the hash.
func HashForRgbaFind(needle []byte, needleWidth int, needleHeight int) uint32 {
  // NOTE: These checks are mainly here to prevent segmentation faults in the
  //       C code. Therefore, panicing is appropriate.
  if len(needle) < needleWidth * needleHeight * 4 {
    panic("Needle width and height do not match buffer size")
  }

  chash := C.GoHashForRgbaFind(unsafe.Pointer(&needle[0]), C.int(needleWidth),
      C.int(needleHeight))
  return uint32(chash)
}

// RgbaFind looks for a needle image in a hastack image.
// It returns the number of matches and the coordinates of the last match.
// The scratch space capacity must be at least 4 * hayWidth. The needle's hash
// can be computed by RgbaHashForFind.
func RgbaFind(haystack []byte, hayWidth int, hayHeight int, needle []byte,
    needleWidth int, needleHeight int, needleHash uint32,
    scratch []byte) (int, int, int) {
  // NOTE: These checks are mainly here to prevent segmentation faults in the
  //       C code. Therefore, panicing is appropriate.
  if len(haystack) < hayWidth * hayHeight * 4 {
    panic("Haystack width and height do not match buffer size")
  }
  if len(needle) < needleWidth * needleHeight * 4 {
    panic("Needle width and height do not match buffer size")
  }
  if cap(scratch) < hayWidth * 4 {
    panic("Insufficent scratch buffer capacity")
  }

  var cmatchX C.int
  var cmatchY C.int
  ccount := C.GoRgbaFind(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(hayHeight),
      C.int(needleWidth), C.int(needleHeight), C.uint32_t(needleHash),
      unsafe.Pointer(&scratch[0]), &cmatchX, &cmatchY)

  return int(ccount), int(cmatchX), int(cmatchY)
}
