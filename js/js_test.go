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
   vis, err := NewVisit(f)
   if err != nil {
      t.Fatal(err)
   }
   for k, v := range vis.Nodes {
      fmt.Println(k, v)
   }
}
