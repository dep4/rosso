package net

import (
   "os"
   "testing"
)

func TestValues(t *testing.T) {
   file, err := os.Open("ignore.txt")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   val := NewValues()
   num, err := val.ReadFrom(file)
   if err != nil {
      t.Fatal(err)
   }
   if num != 679 {
      t.Fatal(err)
   }
   val.WriteTo(os.Stdout)
}
