package hls

import (
   "fmt"
   "os"
   "testing"
)

var rawIVs = []string{
   "00000000000000000000000000000001",
   "0X00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, rawIV := range rawIVs {
      iv, err := Segment{RawIV: &rawIV}.IV()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(iv)
   }
}

func TestSegment(t *testing.T) {
   file, err := os.Open("ignore/apple-segment.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   segs, err := NewScanner(file).Segments()
   if err != nil {
      t.Fatal(err)
   }
   for _, seg := range segs.Key() {
      fmt.Printf("%+v\n", seg)
   }
}
