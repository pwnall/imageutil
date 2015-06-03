package imageutil

import (
  "bytes"
  "crypto/sha256"
  "encoding/hex"
  "reflect"
  "sort"
  "testing"
)

type Pillars [][4]int32
func (p Pillars) Len() int { return len(p) }
func (p Pillars) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p Pillars) Less(i, j int) bool {
  for k := 0; k < 4; k += 1 {
    if p[i][k] < p[j][k] {
      return true
    }
    if p[i][0] > p[j][0] {
      return false
    }
  }
  return false
}

func TestRgbaFindPillars(t *testing.T) {
  goldPillars := [][4]int32{
    {167, 499, 0, 166},
    {169, 501, 0, 168},
    {169, 502, 0, 168},
    {170, 503, 0, 169},
    {171, 504, 0, 170},
    {171, 505, 0, 170},
    {172, 506, 0, 171},
    {173, 507, 0, 172},
    {173, 508, 0, 172},
    {175, 510, 0, 174},
  }


  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  pillars := make([][4]int32, 10)
  RgbaFindPillars(image.Pix, image.Bounds().Dx(), image.Bounds().Dy(),
      230, 255, 150, 220, 0, 120, pillars)

  sort.Sort(Pillars(pillars))
  if !reflect.DeepEqual(goldPillars, pillars) {
    t.Errorf("Incorrect pillars: %v\n", pillars)
  }
}

func TestRgbaResetPuddles(t *testing.T) {
  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }

  changedPixels := make([]byte, len(image.Pix))
  copy(changedPixels, image.Pix)
  RgbaThreshold(changedPixels, 0, 0, 0, 0, 0, 0);

  RgbaResetPuddles(changedPixels)
  if !bytes.Equal(changedPixels, image.Pix) {
    t.Error("Alpha channel not completely restored")
  }
}

func TestRgbaFindPuddles(t *testing.T) {
  goldHash :=
      "5c875aac2f1776d77e9ffe81bf8ff45462521e0e120d7c619a3a484f88eb8536"
  goldHash2 :=
      "c2eb55f067bc942ad0f59c77de46f0021f38d4e69266ed60ec387223854198fd"

  image, err := ReadRgbaPng("test_data/fruits.png")
  if err != nil {
    t.Fatal(err)
  }
  width, height := image.Bounds().Dx(), image.Bounds().Dy()

  puddlePixels := make([][2]int32, 1024 * 768)
  puddleSize := RgbaFindPuddle(image.Pix, width, height,
      230, 255, 150, 220, 0, 120, 0, puddlePixels)

  // Save the "visited" marks for debugging.
  RgbaToPng(image.Pix, image.Bounds().Dx(), image.Bounds().Dy(),
      "test_tmp/fruits_RgbaFindPuddle.png")

  if puddleSize != 17788 {
    t.Error("Incorrect puddle size: ", puddleSize)
  }
  for i := 0; i < puddleSize; i += 1 {
    x, y := int(puddlePixels[i][0]), int(puddlePixels[i][1])
    if x >= width || y >= height {
      t.Errorf("Puddle pixel out of bounds: %d, %d", x, y)
    } else if image.Pix[4 * (y * width + x) + 3] != 0 {
      t.Errorf("Alpha not zero in puddle pixel: %d, %d", x, y)
    }
  }

  hash := sha256.Sum256(image.Pix)
  hexHash := hex.EncodeToString(hash[:])
  if hexHash != goldHash {
    t.Error("Puddle pixel data hash mismatch. Got :", hexHash)
  }


  RgbaResetPuddles(image.Pix)
  puddlePixels = puddlePixels[:1024]  // Limit the puddle size.
  puddleSize = RgbaFindPuddle(image.Pix, width, height,
      230, 255, 150, 220, 0, 120, 100, puddlePixels)

  // Save the "visited" marks for debugging.
  RgbaToPng(image.Pix, image.Bounds().Dx(), image.Bounds().Dy(),
      "test_tmp/fruits_RgbaFindPuddle2.png")

  if puddleSize != 1024 {
    t.Error("Incorrect puddle size: ", puddleSize)
  }
  for i := 0; i < puddleSize; i += 1 {
    x, y := int(puddlePixels[i][0]), int(puddlePixels[i][1])
    if x >= width || y >= height {
      t.Errorf("Puddle 2 pixel out of bounds: %d, %d", x, y)
    } else if image.Pix[4 * (y * width + x) + 3] != 0 {
      t.Errorf("Alpha not zero in puddle 2 pixel: %d, %d", x, y)
    }
  }

  hash = sha256.Sum256(image.Pix)
  hexHash = hex.EncodeToString(hash[:])
  if hexHash != goldHash2 {
    t.Error("Puddle 2 pixel data hash mismatch. Got :", hexHash)
  }
}
