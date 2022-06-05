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
   seg, err := NewScanner(file).Segment(&url.URL{})
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(seg.Key)
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
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
   master, err := NewScanner(file).Master(base)
   if err != nil {
      t.Fatal(err)
   }
   streams := master.Streams.
      Codecs("hvc1").
      Codecs("mp4a").
      RawQuery("cdn=vod-ak-aoc.tv.apple.com").
      VideoRange("PQ")
   for _, stream := range streams {
      fmt.Println(stream)
   }
   stream := master.Streams.GetBandwidth(0)
   fmt.Printf("%a\n", stream)
}

func TestMedia(t *testing.T) {
   file, err := os.Open("m3u8/ignore.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   base, err := url.Parse(playlist)
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
