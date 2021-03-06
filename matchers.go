package imageutil

// #include "c/matchers.c"
import "C"  // cgo

import (
  "unsafe"
)

// RgbaCheckCrop returns true if an image is a cropped version of another one.
// This is image pattern-matching, but only aligns the pattern with the image
// in one predetermined position.
func RgbaCheckCrop(haystack []byte, hayWidth int, hayHeight int,
    needle []byte, needleWidth int, needleHeight int, needleLeft int,
    needleTop int) bool {

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
  if needleLeft < 0 || needleLeft + needleWidth > hayWidth {
    return false
  }
  if needleTop < 0 || needleTop + needleHeight > hayHeight {
    return false
  }

  // NOTE: The haystack's height is irrelevant to the actual matching logic,
  //       so it is omitted.
  cresult := C.GoRgbaCheckCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(needleWidth),
      C.int(needleHeight), C.int(needleLeft), C.int(needleTop))
  return cresult != 0
}

// RgbaCheckMaskedCrop checks if an image is a crop&mask from another image.
// This is image pattern-matching, but only aligns the pattern with the image
// in one predetermined position, and masks the image on the fly. The pattern
// is assumed to have already been masked.
func RgbaCheckMaskedCrop(haystack []byte, hayWidth int, hayHeight int,
    needle []byte, needleWidth int, needleHeight int, needleLeft int,
    needleTop int, rgbaMask uint32) bool {

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
  if needleLeft < 0 || needleLeft + needleWidth > hayWidth {
    return false
  }
  if needleTop < 0 || needleTop + needleHeight > hayHeight {
    return false
  }

  // RGBA -> ARGB, because Intel is little-endian.
  argbMask := uint32(((rgbaMask & 0xff) << 24) | ((rgbaMask & 0xff00) << 8) |
      ((rgbaMask & 0xff0000) >> 8) | ((rgbaMask & 0xff000000) >> 24))

  // NOTE: The haystack's height is irrelevant to the actual matching logic,
  //       so it is omitted.
  cresult := C.GoRgbaCheckMaskedCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(needleWidth),
      C.int(needleHeight), C.int(needleLeft), C.int(needleTop),
      C.uint32_t(argbMask))
  return cresult != 0
}

// RgbaDiffMaskedCrop diffs an image with a crop&mask of another image.
// It returns the sum of absolute pixel differences.
func RgbaDiffMaskedCrop(haystack []byte, hayWidth int, hayHeight int,
    needle []byte, needleWidth int, needleHeight int, needleLeft int,
    needleTop int, rgbaMask uint32) int64 {

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
  if needleLeft < 0 || needleLeft + needleWidth > hayWidth {
    return 0
  }
  if needleTop < 0 || needleTop + needleHeight > hayHeight {
    return 0
  }

  // RGBA -> ARGB, because Intel is little-endian.
  argbMask := uint32(((rgbaMask & 0xff) << 24) | ((rgbaMask & 0xff00) << 8) |
      ((rgbaMask & 0xff0000) >> 8) | ((rgbaMask & 0xff000000) >> 24))

  // NOTE: The haystack's height is irrelevant to the actual matching logic,
  //       so it is omitted.
  cresult := C.GoRgbaDiffMaskedCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(needleWidth),
      C.int(needleHeight), C.int(needleLeft), C.int(needleTop),
      C.uint32_t(argbMask))
  return int64(cresult)
}

// RgbaDiffThresholdCrop diffs an image with a crop&threshold of another image.
// It returns the sum of absolute pixel differences.
func RgbaDiffThresholdCrop(haystack []byte, hayWidth int, hayHeight int,
    needle []byte, needleWidth int, needleHeight int, needleLeft int,
    needleTop int, minRed int, maxRed int, minGreen int, maxGreen int,
    minBlue int, maxBlue int) int {
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
  if needleLeft < 0 || needleLeft + needleWidth > hayWidth {
    return 0
  }
  if needleTop < 0 || needleTop + needleHeight > hayHeight {
    return 0
  }

  // NOTE: The haystack's height is irrelevant to the actual matching logic,
  //       so it is omitted.
  cresult := C.GoRgbaDiffThresholdCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(needleWidth),
      C.int(needleHeight), C.int(needleLeft), C.int(needleTop),
      C.uint8_t(minRed), C.uint8_t(minGreen), C.uint8_t(minBlue),
      C.uint8_t(maxRed), C.uint8_t(maxGreen), C.uint8_t(maxBlue))
  return int(cresult)
}


// HashForRgbaFindCrop computes the needle hash needed by RgbaFind.
// It returns the hash.
func HashForRgbaFindCrop(needle []byte, needleWidth int,
    needleHeight int) uint32 {
  // NOTE: These checks are mainly here to prevent segmentation faults in the
  //       C code. Therefore, panicing is appropriate.
  if len(needle) < needleWidth * needleHeight * 4 {
    panic("Needle width and height do not match buffer size")
  }

  chash := C.GoHashForRgbaFindCrop(unsafe.Pointer(&needle[0]),
      C.int(needleWidth), C.int(needleHeight))
  return uint32(chash)
}

// RgbaFindCrop looks for a needle image in a hastack image.
// It returns the number of matches and the coordinates of the last match.
// The scratch space capacity must be at least 4 * hayWidth. The needle's hash
// can be computed by RgbaHashForFindCrop.
func RgbaFindCrop(haystack []byte, hayWidth int, hayHeight int, needle []byte,
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

  var cmatchLeft C.int
  var cmatchTop C.int
  ccount := C.GoRgbaFindCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(hayHeight),
      C.int(needleWidth), C.int(needleHeight), C.uint32_t(needleHash),
      unsafe.Pointer(&scratch[0]), &cmatchLeft, &cmatchTop)

  return int(ccount), int(cmatchLeft), int(cmatchTop)
}

// RgbaFindMaskedCrop looks for a masked needle image in a hastack image.
// It returns the number of matches and the coordinates of the last match.
// The scratch space capacity must be at least 4 * hayWidth. The needle's hash
// can be computed by RgbaHashForFindCrop. The needle is assumed to have been
// masked before RgbaHashForFindCrop and this method are called.
func RgbaFindMaskedCrop(haystack []byte, hayWidth int, hayHeight int,
    needle []byte, needleWidth int, needleHeight int, rgbaMask uint32,
    needleHash uint32, scratch []byte) (int, int, int) {
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

  // RGBA -> ARGB, because Intel is little-endian.
  argbMask := uint32(((rgbaMask & 0xff) << 24) | ((rgbaMask & 0xff00) << 8) |
      ((rgbaMask & 0xff0000) >> 8) | ((rgbaMask & 0xff000000) >> 24))

  var cmatchLeft C.int
  var cmatchTop C.int
  ccount := C.GoRgbaFindMaskedCrop(unsafe.Pointer(&haystack[0]),
      unsafe.Pointer(&needle[0]), C.int(hayWidth), C.int(hayHeight),
      C.int(needleWidth), C.int(needleHeight), C.uint32_t(argbMask),
      C.uint32_t(needleHash), unsafe.Pointer(&scratch[0]), &cmatchLeft,
      &cmatchTop)

  return int(ccount), int(cmatchLeft), int(cmatchTop)
}
