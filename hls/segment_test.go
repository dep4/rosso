package hls

import (
   "fmt"
   "net/http"
   "net/url"
   "os"
   "testing"
)

var ivs = []string{
   "00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, iv := range ivs {
      buf, err := hexDecode(iv)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(buf)
   }
}

func TestSegment(t *testing.T) {
   seg, err := newSegment()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(seg.Key.URI.String())
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   block, err := NewCipher(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("ignore.ts")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for _, info := range seg.Info {
      fmt.Println(info)
      res, err := http.Get(info.URI.String())
      if err != nil {
         t.Fatal(err)
      }
      defer res.Body.Close()
      if _, err := block.Copy(file, res.Body, info.IV); err != nil {
         t.Fatal(err)
      }
   }
}

func newSegment() (*Segment, error) {
   file, err := os.Open("m3u8/abc-segment.m3u8")
   if err != nil {
      return nil, err
   }
   defer file.Close()
   return NewSegment(&url.URL{}, file)
}
