package hls

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "net/http"
   "os"
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

const segmentURL = 
   "https://ga.video.cdn.pbs.org/videos/nature" +
   "/1803a543-5d57-41f6-81bf-d199448d45f6/2000281849/hd-16x9-mezzanine-1080p" +
   "/naat4008_r-hls-16x9-1080pAudio+Selector+1.m3u8"

func TestSegment(t *testing.T) {
   res, err := http.Get(segmentURL)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   seg, err := NewScanner(res.Body).Segment(res.Request.URL)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.aac")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   pro := format.ProgressChunks(file, len(seg.Info))
   for _, info := range seg.Info {
      res, err := http.Get(info.URI.String())
      if err != nil {
         t.Fatal(err)
      }
      pro.AddChunk(res.ContentLength)
      if _, err := io.Copy(pro, res.Body); err != nil {
         t.Fatal(err)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}
