package hls

import (
   "fmt"
   "sort"
   "testing"
)

func TestMaster(t *testing.T) {
   master, err := newMaster()
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Stream {
      fmt.Println(stream)
   }
}

func TestProgress(t *testing.T) {
   seg := Segment{
      Info: make([]Information, 9),
   }
   for i := range seg.Info {
      fmt.Print(seg.Progress(i))
   }
   fmt.Println("END")
}

func TestSort(t *testing.T) {
   mas := &Master{Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 144},
      {Bandwidth: 1080},
      {Bandwidth: 720},
      {Bandwidth: 2160},
   }}
   sort.Sort(Bandwidth{mas, 720})
   for _, str := range mas.Stream {
      fmt.Println(str)
   }
}
