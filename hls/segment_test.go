package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func TestSegment(t *testing.T) {
   file, err := os.Open("m3u8/abc-segment.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   seg, err := NewSegment(&url.URL{}, file)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", seg)
}
