package os

import (
   "testing"
)

func Test_Create(t *testing.T) {
   file, err := Create("ignore.txt")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
}
