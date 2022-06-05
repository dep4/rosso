package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func TestMedia(t *testing.T) {
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
   media := master.Media.
      GroupID("stereo").
      Name("English").
      RawQuery("cdn=vod-ak-aoc.tv.apple.com").
      Type("AUDIO")
   for _, medium := range media {
      fmt.Println(medium)
   }
}
