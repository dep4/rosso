package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

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

func TestCwtvMaster(t *testing.T) {
   file, err := os.Open("m3u8/master-cwtv.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   mas, err := NewMaster(&url.URL{}, file)
   if err != nil {
      t.Fatal(err)
   }
   for _, med := range mas.Media {
      fmt.Printf("%+v\n", med)
   }
   for _, str := range mas.Stream {
      fmt.Println(str)
   }
}
