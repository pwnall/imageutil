package imageutil

// This file doesn't have any test cases. It only contains TestMain, which has
// setup and teardown code for all the tests.

import (
  "fmt"
  "os"
  "testing"
)

func TestMain(m *testing.M) {
  // The test images will be stored in test_tmp.
  if err := os.MkdirAll("test_tmp", 0777); err != nil {
    fmt.Printf("%v\n", err)
    os.Exit(1)
  }

  result := m.Run()
  os.Exit(result)
}
