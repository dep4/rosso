package js

import (
   "fmt"
   "os"
   "testing"
)

func TestJS(t *testing.T) {
   f, err := os.Open("ig.js")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   p, err := Parse(f)
   if err != nil {
      t.Fatal(err)
   }
   for k, v := range p {
      fmt.Println(k, v)
   }
}
