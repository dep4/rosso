package hls

import (
   "fmt"
   "os"
   "testing"
)

func TestSegment(t *testing.T) {
   file, err := os.Open("ignore/apple-segment.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   seg, err := NewScanner(file).Segment()
   if err != nil {
      t.Fatal(err)
   }
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
   key, err := seg.Base64()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(key)
}
