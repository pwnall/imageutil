package imageutil

import (
  "bytes"
  "encoding/hex"
  "crypto/sha256"
  "testing"
)

func TestReadRgbaPng(t *testing.T) {
  goldImageHash :=
      "d333fb91aa05709483df2d00f62e3caa91db0be1b30ff72e8e829f3264cb30b9"

  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  if image.Rect.Min.X != 0 || image.Rect.Max.X != 512 ||
      image.Rect.Min.Y != 0 || image.Rect.Max.Y != 512 {
    t.Errorf("Incorrect image rectangle: %v\n", image.Rect)
  }

  hash := sha256.Sum256(image.Pix)
  imageHash := hex.EncodeToString(hash[:])
  if imageHash != goldImageHash {
    t.Error("Pixel data hash mismatch. Got :", imageHash)
  }
}

func TestRgbaToPng(t *testing.T) {
  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  err = RgbaToPng(image.Pix, image.Bounds().Dx(), image.Bounds().Dy(),
      "test_tmp/fruits_RgbaToPng.png")
  if err != nil {
    t.Fatal(err)
  }

  recodedImage, err := ReadRgbaPng("test_tmp/fruits_RgbaToPng.png")
  if err != nil {
    t.Fatal(err)
  }
  if !bytes.Equal(image.Pix, recodedImage.Pix) {
    t.Error("Pixel data mismatch")
  }
}
