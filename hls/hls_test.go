package hls

import (
   "fmt"
   "github.com/89z/format"
   "io"
   "net/http"
   "net/url"
   "os"
   "sort"
   "testing"
)

var ivs = []string{
   "00000000000000000000000000000001",
   "0X00000000000000000000000000000001",
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

var cbc = &url.URL{
   Scheme: "https",
   Host: "cbcrcott-gem.akamaized.net",
   Path: "/0f73fb9d-87f0-4577-81d1-e6e970b89a69/CBC_DOWNTON_ABBEY_S01E05.ism/desktop_master.m3u8",
}

func TestSegment(t *testing.T) {
   file, err := os.Open("m3u8/cbc-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := NewScanner(file).Master(cbc)
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Stream {
      if stream.Bandwidth == 621388 {
         res, err := http.Get(stream.URI.String())
         if err != nil {
            t.Fatal(err)
         }
         defer res.Body.Close()
         seg, err := NewScanner(res.Body).Segment(res.Request.URL)
         if err != nil {
            t.Fatal(err)
         }
         file, err := os.Create("ignore.mp4")
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
   }
}

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
