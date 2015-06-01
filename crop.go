package imageutil

import (
  "image"
  "image/draw"
)

// CropRgba crops an RGBA image into a target slice.
// The target slice's length is set to the needed image length. If the slice's
// capacity is too small, the slice is re-created.
func CropRgba(rawImage []byte, width int, height int, xOffset int,
    yOffset int, xSize int, ySize int, target *[]byte) {

  source := image.RGBA{Pix: rawImage, Stride: width * 4,
      Rect: image.Rect(0, 0, width, height)}

  targetSize := xSize * ySize * 4
  if cap(*target) < targetSize {
    *target = make([]byte, targetSize, targetSize)
  } else if len(*target) != targetSize {
    *target = (*target)[:targetSize]
  }
  cropped := image.RGBA{Pix: *target, Stride: xSize * 4,
      Rect: image.Rect(0, 0, xSize, ySize)}

  sourcePt := image.Pt(xOffset, yOffset)
  draw.Draw(&cropped, cropped.Bounds(), &source, sourcePt, draw.Src)
}
