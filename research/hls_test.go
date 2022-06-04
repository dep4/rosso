package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const base = "https://play.itunes.apple.com" +
   "/WebObjects/MZPlay.woa/hls/subscription/playlist.m3u8"

func TestMaster(t *testing.T) {
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
   streams := master.Streams.Streams(func(s Stream) bool {
      return s.URI.Query().Get("cdn") == "vod-ak-aoc.tv.apple.com"
   })
   for _, stream := range streams {
      fmt.Println(stream)
   }
}

/*
Codecs:ac-3,dvh1.05.06
Codecs:ac-3,dvh1.05.06
Codecs:ac-3,dvh1.05.06
Codecs:ac-3,hvc1.2.20000000.H150.B0
Codecs:ac-3,hvc1.2.20000000.H150.B0
Codecs:ac-3,hvc1.2.20000000.H150.B0

Codecs:dvh1.05.06,ec-3
Codecs:dvh1.05.06,ec-3
Codecs:dvh1.05.06,ec-3
Codecs:ec-3,hvc1.2.20000000.H150.B0
Codecs:ec-3,hvc1.2.20000000.H150.B0
Codecs:ec-3,hvc1.2.20000000.H150.B0

Codecs:dvh1.05.06,mp4a.40.2
Codecs:dvh1.05.06,mp4a.40.2
Codecs:dvh1.05.06,mp4a.40.2
Codecs:hvc1.2.20000000.H150.B0,mp4a.40.2
Codecs:hvc1.2.20000000.H150.B0,mp4a.40.2
Codecs:hvc1.2.20000000.H150.B0,mp4a.40.2
*/
