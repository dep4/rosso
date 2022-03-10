package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

var masters = []Master{
   {Stream: []Stream{
      {Bandwidth: 144},
      {Bandwidth: 480},
      {Bandwidth: 720},
      {Bandwidth: 1080},
      {Bandwidth: 2160},
   }},
   {Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 1080},
   }},
   {Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 2160},
   }},
   {Stream: []Stream{
      {Bandwidth: 144},
      {Bandwidth: 1080},
   }},
   {Stream: []Stream{
      {Bandwidth: 480},
   }},
   {Stream: []Stream{
   }},
}

func TestQuality(t *testing.T) {
   for _, master := range masters {
      fmt.Printf("%#v\n", master.GetStream(720))
   }
   fmt.Printf("%#v\n", masters[0].GetStream(-1))
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
