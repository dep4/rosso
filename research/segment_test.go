package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const playlist = "https://play.itunes.apple.com" +
   "/WebObjects/MZPlay.woa/hls/subscription/playlist.m3u8"

func TestSegment(t *testing.T) {
   file, err := os.Open("ignore/apple-segment.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   seg, err := NewScanner(file).Segment()
   if err != nil {
      t.Fatal(err)
   }
   base, err := url.Parse(playlist)
   if err != nil {
      t.Fatal(err)
   }
   key, err := seg.Key(base)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(key)
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
}
