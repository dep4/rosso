package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

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
