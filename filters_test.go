package imageutil

import (
  "encoding/hex"
  "crypto/sha256"
  "testing"
)

func TestBuildMask(t *testing.T) {
  cases := map[uint32]uint64{
    0x12345678: 0x7856341278563412,
    0xff000000: 0x000000ff000000ff,
    0x00ff0000: 0x0000ff000000ff00,
    0x0000ff00: 0x00ff000000ff0000,
    0x000000ff: 0xff000000ff000000,
  }

  for input, golden := range cases {
    if output := BuildRgbaMask(input); output != golden {
      t.Errorf("Got unexpected value %X for input %X\n", output, input)
    }
  }
}

func TestMaskRgba(t *testing.T) {
  goldHash :=
      "bef35a27034ada1a7f35c21e82ba0bb59f8a727f9d6ce04339cb13d800d8570a"

  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  MaskRgba(image.Pix, BuildRgbaMask(0xf0e0c0ff))
  // Save the crop result for debugging.
  RgbaToPng(image.Pix, image.Bounds().Dx(), image.Bounds().Dy(),
      "test_tmp/fruits_RgbaMask.png")

  hash := sha256.Sum256(image.Pix)
  hexHash := hex.EncodeToString(hash[:])
  if hexHash != goldHash {
    t.Error("Pixel data hash mismatch. Got :", hexHash)
  }
}
