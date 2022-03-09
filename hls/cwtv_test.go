package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func TestRemove(t *testing.T) {
   str := Stream{
      Bandwidth: 1, Codecs: "Codecs",
      URI: &url.URL{Scheme: "http", Host: "example.com"},
   }
   addr := str.RemoveURI()
   fmt.Println(addr)
   fmt.Println(str)
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
   for _, str := range mas.Stream {
      fmt.Println(str)
   }
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
