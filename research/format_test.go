package format

import (
   "github.com/89z/format/hls"
   "io"
   "net/http"
   "os"
   "testing"
)

func TestProgress(t *testing.T) {
   pro := NewProgress(io.Discard, 1)
   res, err := http.Get("https://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   if _, err := pro.Copy(res); err != nil {
      t.Fatal(err)
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
   seg, err := hls.NewScanner(res.Body).Segment(res.Request.URL)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.aac")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   pro := NewProgress(file, len(seg.Info))
   for _, info := range seg.Info {
      res, err := http.Get(info.URI.String())
      if err != nil {
         t.Fatal(err)
      }
      if _, err := pro.Copy(res); err != nil {
         t.Fatal(err)
      }
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}
