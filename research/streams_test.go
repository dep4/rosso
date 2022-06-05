package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const base = "https://play.itunes.apple.com" +
   "/WebObjects/MZPlay.woa/hls/subscription/playlist.m3u8"

func TestStream(t *testing.T) {
   file, err := os.Open("ignore.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   base, err := url.Parse(base)
   if err != nil {
      t.Fatal(err)
   }
   master, err := NewScanner(file).Master(base)
   if err != nil {
      t.Fatal(err)
   }
   streams := master.Streams.
      Codec("hvc1").
      Codec("mp4a").
      RawQuery("cdn=vod-ak-aoc.tv.apple.com").
      VideoRange("PQ")
   for _, stream := range streams {
      fmt.Println(stream)
   }
}
