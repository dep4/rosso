package hls

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "testing"
)

const base = "https://play.itunes.apple.com" +
   "/WebObjects/MZPlay.woa/hls/subscription/playlist.m3u8"

func TestSegment(t *testing.T) {
   file, err := os.Open("m3u8/cbc-video.m3u8")
   if err != nil {
      return nil, err
   }
   defer file.Close()
   addr, err := url.Parse(cbcSegment)
   if err != nil {
      return nil, err
   }
   return NewScanner(file).Segment(addr)
   // FIXME
   fmt.Println("GET", seg.Key)
   res, err := http.Get(seg.Key.String())
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   block, err := NewCipher(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for _, info := range seg.Info {
      fmt.Println("GET", info.URI)
      res, err := http.Get(info.URI.String())
      if err != nil {
         t.Fatal(err)
      }
      if res.StatusCode != http.StatusOK {
         t.Fatal(res.Status)
      }
      if _, err := block.Copy(file, res.Body, info.IV); err != nil {
         t.Fatal(err)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}

func TestStream(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
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
