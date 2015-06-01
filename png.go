package imageutil

import (
  "errors"
  "image"
  "image/png"
  "os"
)

// RgbaToPng encodes a raw RGBA-encoded image into a PNG image.
// It returns any error encountered.
func RgbaToPng(rawImage []byte, width int, height int, fileName string) error {
  // NOTE: This hack wraps an RGBA structure over an existing slice, to avoid
  //       a memory copy.
  rgbaImage := image.RGBA{Pix: rawImage, Stride: (width * 4),
    Rect: image.Rect(0, 0, width, height)}

  f, err := os.Create(fileName)
  if err != nil {
    return err
  }
  defer f.Close()

  return png.Encode(f, &rgbaImage)
}

// ReadRgbaPng decodes a PNG image from a file into a raw RGBA buffer.
// It returns the decoded image and any error encountered.
func ReadRgbaPng(fileName string) (*image.RGBA, error) {
  f, err := os.Open(fileName)
  if err != nil {
    return nil, err
  }
  defer f.Close()

  pngImage, err := png.Decode(f)
  if err != nil {
    return nil, err
  }

  rgbaImage, success := pngImage.(*image.RGBA)
  if !success {
    return nil, errors.New("Decoded image is not RGBA")
  }
  return rgbaImage, nil
}
