package imageutil

import (
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
      230, 255, 150, 220, 0, 120, 10, pillars)

  sort.Sort(Pillars(pillars))
  if !reflect.DeepEqual(goldPillars, pillars) {
    t.Errorf("Incorrect pillars: %v\n", pillars)
  }
}

