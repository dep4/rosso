package js

import (
   "os"
   "testing"
)

func TestJS(t *testing.T) {
   f, err := os.Open("ig.js")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   if _, err := statements(f); err != nil {
      t.Fatal(err)
   }
}
