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
      iv, err := Information{RawIV: rawIV}.IV()
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
   seg, err := NewScanner(file).Segment()
   if err != nil {
      t.Fatal(err)
   }
   for _, info := range seg.Info {
      fmt.Printf("%+v\n", info)
   }
   pssh, err := seg.PSSH()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(pssh)
}
