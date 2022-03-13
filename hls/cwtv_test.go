package hls

import (
   "fmt"
   "net/url"
   "os"
   "sort"
   "testing"
)

func TestCwtvMaster(t *testing.T) {
   file, err := os.Open("m3u8/master-cwtv.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   mas, err := NewMaster(&url.URL{}, file)
   if err != nil {
      t.Fatal(err)
   }
   for _, str := range mas.Stream {
      fmt.Println(str)
      // fmt.Printf("%u\n", str)
   }
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

func TestProgress(t *testing.T) {
   seg := Segment{
      Info: make([]Information, 9),
   }
   for i := range seg.Info {
      fmt.Print(seg.Progress(i))
   }
   fmt.Println("END")
}

func TestCwtvSegment(t *testing.T) {
   file, err := os.Open("m3u8/segment-cwtv.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   seg, err := NewSegment(&url.URL{}, file)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(seg)
}
