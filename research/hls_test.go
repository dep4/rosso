package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const playlist = "https://play.itunes.apple.com" +
   "/WebObjects/MZPlay.woa/hls/subscription/playlist.m3u8"

func TestMedia(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   NewScanner(file).Master(&url.URL{})
}

func TestStream(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   base, err := url.Parse(playlist)
   if err != nil {
      t.Fatal(err)
   }
   master, err := NewScanner(file).Master2(base)
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Printf("%a\n", stream)
   }
}
