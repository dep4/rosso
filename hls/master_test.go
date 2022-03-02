package hls

import (
   "fmt"
   "testing"
)

func TestMaster(t *testing.T) {
   mas, err := doMaster()
   if err != nil {
      t.Fatal(err)
   }
   for _, med := range mas.Media {
      fmt.Printf("%+v\n", med)
   }
   for _, str := range mas.Stream {
      fmt.Printf("%+v\n", str)
   }
}
