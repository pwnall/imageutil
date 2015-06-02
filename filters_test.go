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

func TestRgbPixelToHsl(t *testing.T) {
  cases := [][6]int {
    {255, 0, 0, 0, 255, 127},
    {0, 255, 0, 84, 255, 127},
    {0, 0, 255, 168, 255, 127},
    {191, 64, 64, 0, 127, 127},
    {64, 191, 64, 84, 127, 127},
    {64, 64, 191, 168, 127, 127},
    {223, 159, 159, 0, 127, 191},
    {159, 223, 159, 84, 127, 191},
    {159, 159, 223, 168, 127, 191},
    {191, 191, 64, 42, 127, 127},
    {64, 191, 191, 126, 127, 127},
    {191, 64, 191, 214, 127, 127},
    {207, 175, 197, 228, 63, 191},
  }

  for _, testCase := range cases {
    r, g, b := testCase[0], testCase[1], testCase[2]
    hsl := [3]int{testCase[3], testCase[4], testCase[5]}

    gotH, gotS, gotL := RgbPixelToHsl(r, g, b)
    gotHsl := [3]int{gotH, gotS, gotL}

    if hsl != gotHsl {
      t.Errorf("Failed on %v, got %v\n", testCase, gotHsl)
    }
  }
}

func TestRgbaToHsla(t *testing.T) {
  goldHash :=
      "bef35a27034ada1a7f35c21e82ba0bb59f8a727f9d6ce04339cb13d800d8570a"

  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  hslaBytes := make([]byte, len(image.Pix))
  RgbaToHsla(image.Pix, hslaBytes)
  // Save the crop result for debugging.
  RgbaToPng(hslaBytes, image.Bounds().Dx(), image.Bounds().Dy(),
      "test_tmp/fruits_RgbaToHsla.png")

  hash := sha256.Sum256(image.Pix)
  hexHash := hex.EncodeToString(hash[:])
  if hexHash != goldHash {
    t.Error("Pixel data hash mismatch. Got :", hexHash)
  }
}
