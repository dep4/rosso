package hls

import (
   "fmt"
   "net/http"
   "sort"
   "testing"
)

func TestBandwidth(t *testing.T) {
   master := &Master{Stream: []Stream{
      {Bandwidth: 480},
      {Bandwidth: 144},
      {Bandwidth: 1080},
      {Bandwidth: 720},
      {Bandwidth: 2160},
   }}
   sort.Sort(Bandwidth{master, 720})
   for _, str := range master.Stream {
      fmt.Println(str)
   }
}

var ivs = []string{
   "00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, iv := range ivs {
      buf, err := scanHex(iv)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(buf)
   }
}

////////////////////////////////////////////////////////////////////////////////

const segmentURL = 
   "https://ga.video.cdn.pbs.org/videos/nature" +
   "/1803a543-5d57-41f6-81bf-d199448d45f6/2000281849/hd-16x9-mezzanine-1080p" +
   "/naat4008_r-hls-16x9-1080pAudio+Selector+1.m3u8"

const masterURL =
   "https://link.theplatform.com/s/dJ5BDC/media/guid/2198311517" +
   "/uOkB1qnYXZXeAM34djVYNHje2_gK4mmO?formats=MPEG4,M3U"

func TestSegment(t *testing.T) {
   fmt.Println("GET", masterURL)
   res, err := http.Get(masterURL)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   master, err := NewScanner(res.Body).Master(res.Request.URL)
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Stream {
      if stream.Bandwidth == 497000 {
         fmt.Println("GET", stream.URI)
         res, err := http.Get(stream.URI.String())
         if err != nil {
            t.Fatal(err)
         }
         defer res.Body.Close()
         seg, err := NewScanner(res.Body).Segment(stream.URI)
         if err != nil {
            t.Fatal(err)
         }
         length := seg.Length(stream)
         if length != 82456897 {
            t.Fatal(length)
         }
      }
   }
}
