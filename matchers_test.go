package imageutil

import (
  "testing"
)

func TestRgbaCheckCrop(t *testing.T) {
  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  width, height := image.Bounds().Dx(), image.Bounds().Dy()
  xOffset, yOffset := 200, 400
  xSize, ySize := 128, 16

  var cropBytes []byte
  CropRgba(image.Pix, width, height, xOffset, yOffset, xSize, ySize,
      &cropBytes)

  result := RgbaCheckCrop(image.Pix, width, height, cropBytes, xSize, ySize,
      xOffset, yOffset)
  if result != true {
    t.Error("Did not detect correctly aligned crop: ", result)
  }

  result = RgbaCheckCrop(image.Pix, width, height, cropBytes, xSize, ySize,
      xOffset - 1, yOffset)
  if result != false {
    t.Error("Did not bounce crop misaligned by (-1, 0): ", result)
  }

  result = RgbaCheckCrop(image.Pix, width, height, cropBytes, xSize, ySize,
      xOffset, yOffset + 1)
  if result != false {
    t.Error("Did not bounce crop misaligned by (0, +1): ", result)
  }
}

func TestRgbaCheckMaksedCrop(t *testing.T) {
  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  width, height := image.Bounds().Dx(), image.Bounds().Dy()
  xOffset, yOffset := 200, 400
  xSize, ySize := 128, 16

  var cropBytes []byte
  CropRgba(image.Pix, width, height, xOffset, yOffset, xSize, ySize,
      &cropBytes)

  result := RgbaCheckMaskedCrop(image.Pix, width, height, cropBytes, xSize,
      ySize, xOffset, yOffset, 0xffffffff)
  if result != true {
    t.Error("Did not detect correctly aligned crop: ", result)
  }

  result = RgbaCheckMaskedCrop(image.Pix, width, height, cropBytes, xSize,
      ySize, xOffset - 1, yOffset, 0xffffffff)
  if result != false {
    t.Error("Did not bounce crop misaligned by (-1, 0): ", result)
  }

  result = RgbaCheckMaskedCrop(image.Pix, width, height, cropBytes, xSize,
      ySize, xOffset, yOffset + 1, 0xffffffff)
  if result != false {
    t.Error("Did not bounce crop misaligned by (0, +1): ", result)
  }

  var maskCropBytes []byte
  CropRgba(image.Pix, width, height, xOffset, yOffset, xSize, ySize,
      &maskCropBytes)
  MaskRgba(maskCropBytes, BuildRgbaMask(0xc0e080ff);
}


func TestRgbaCheckCrop_positives(t *testing.T) {
  cases := [][4]int {
    { 0, 0, 16, 8 },
    { 0, 10, 16, 8 },
    { 10, 0, 16, 8 },
    { 10, 10, 16, 8 },
  }

  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }
  width, height := image.Bounds().Dx(), image.Bounds().Dy()

  var cropBytes []byte
  for _, testCase := range cases {
    xOffset, yOffset := testCase[0], testCase[1]
    xSize, ySize := testCase[2], testCase[3]

    CropRgba(image.Pix, width, height, xOffset, yOffset, xSize, ySize,
        &cropBytes)

    if !RgbaCheckCrop(image.Pix, width, height, cropBytes, xSize, ySize,
        xOffset, yOffset) {
      t.Errorf("Did not detect crop %v\n", testCase)
    }
  }
}


func TestRgbaFind(t *testing.T) {
  cases := [][4]int {
    { 0, 0, 16, 8 },
    { 0, 10, 16, 8 },
    { 10, 0, 16, 8 },
    { 10, 10, 16, 8 },
    { 0, 0, 12, 84 },
    { 500, 0, 12, 8 },
    { 0, 300, 12, 84 },
    { 500, 300, 12, 84 },
  }

  // NOTE: We crop the initial image because we want different width / height,
  //       to make sure that the Rabin-Karp implementation uses width / height
  //       correctly.

  originalImage, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  var imageBytes []byte
  width := 512
  height := 384
  CropRgba(originalImage.Pix, originalImage.Bounds().Dx(),
      originalImage.Bounds().Dy(), originalImage.Bounds().Dx() - width,
      originalImage.Bounds().Dy() - height, width, height, &imageBytes)

  var cropBytes []byte
  scratch := make([]byte, width * 4)
  for _, testCase := range cases {
    xOffset, yOffset := testCase[0], testCase[1]
    xSize, ySize := testCase[2], testCase[3]

    CropRgba(imageBytes, width, height, xOffset, yOffset, xSize, ySize,
        &cropBytes)

    if !RgbaCheckCrop(imageBytes, width, height, cropBytes, xSize, ySize,
        xOffset, yOffset) {
      t.Fatal("RgbaCheckCrop doesn't verify golden crop")
    }

    hash := HashForRgbaFind(cropBytes, xSize, ySize)
    t.Logf("Needle hash: %d\n", hash)
    count, matchX, matchY := RgbaFind(imageBytes, width, height, cropBytes,
        xSize, ySize, hash, scratch)
    if count != 1 || matchX != xOffset || matchY != yOffset {
      t.Errorf("Wrong answer on case %v - count %d, matchX %d, matchY %d",
        testCase, count, matchX, matchY)
    }
  }
}
