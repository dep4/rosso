package hls

import (
   "fmt"
   "net/url"
   "os"
   "strings"
   "testing"
)

const base = "https://play.itunes.apple.com" +
   "/WebObjects/MZPlay.woa/hls/subscription/playlist.m3u8"

func TestMaster(t *testing.T) {
   fn := func(s Stream) bool {
      if s.URI.Query().Get("cdn") != "vod-ak-aoc.tv.apple.com" {
         return false
      }
      if !strings.Contains(s.Codecs, "mp4a") {
         return false
      }
      if !strings.Contains(s.Codecs, "hvc1") {
         return false
      }
      if s.VideoRange != "PQ" {
         return false
      }
      return true
   }
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
   streams := master.Streams.Streams(fn)
   for _, stream := range streams {
      fmt.Println(stream)
   }
}
