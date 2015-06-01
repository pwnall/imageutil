package imageutil

import (
  "bytes"
  "encoding/hex"
  "crypto/sha256"
  "testing"
)

func TestCropRgba(t *testing.T) {
  goldHash :=
      "4725a80241458c161db589a64ec7efb378a4e449a70924c0a2a3043408e92495"

  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  width, height := image.Bounds().Dx(), image.Bounds().Dy()
  xOffset, yOffset := 200, 400
  xSize, ySize := 128, 16

  // Test the case where the slice has to be re-allocated.
  var target []byte
  CropRgba(image.Pix, width, height, xOffset, yOffset, xSize, ySize, &target)
  // Save the crop result for debugging.
  RgbaToPng(target, xSize, ySize, "test_tmp/fruits_Crop.png")
  if len(target) != 8192 {
    t.Fatal("Incorrect crop data size: ", len(target))
  }

  hash := sha256.Sum256(target)
  hexHash := hex.EncodeToString(hash[:])
  if hexHash != goldHash {
    t.Error("Crop pixel data hash mismatch. Got :", hexHash)
  }

  // Test the case where the slice's length has to change.
  target2 := make([]byte, 65536)
  CropRgba(image.Pix, width, height, xOffset, yOffset, xSize, ySize, &target2)
  // Save the crop result for debugging.
  RgbaToPng(target, xSize, ySize, "test_tmp/fruits_Crop.png")
  if len(target2) != 8192 {
    t.Fatal("Incorrect crop2 data size: ", len(target))
  }

  if !bytes.Equal(target2, target) {
    t.Error("Crop2 pixel data mismatch")
  }
}
