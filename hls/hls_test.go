package hls

import (
   "fmt"
   "net/url"
   "os"
   "sort"
   "testing"
)

func TestSegment(t *testing.T) {
   file, err := os.Open("m3u8/paramount-segment.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   seg, err := NewScanner(file).Segment(&url.URL{})
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", seg.Key)
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
}

func TestBandwidth(t *testing.T) {
   master := &Master{Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 144},
      {Bandwidth: 1080},
      {Bandwidth: 720},
      {Bandwidth: 2160},
   }}
   sort.Sort(Bandwidth{master, 720})
   for _, str := range master.Stream {
      fmt.Println(str)
   }
}

var ivs = []string{
   "00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, iv := range ivs {
      buf, err := scanHex(iv)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(buf)
   }
}
