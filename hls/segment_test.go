package hls

import (
   "fmt"
   "os"
   "testing"
)

var names = []string{
   "m3u8/apple-audio.m3u8",
   "m3u8/cbc-video.m3u8",
   "m3u8/roku-segment.m3u8",
}

func Test_Segment(t *testing.T) {
   for _, name := range names {
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      seg, err := New_Scanner(file).Segment()
      if err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n\n", seg)
   }
}

var raw_ivs = []string{
   "00000000000000000000000000000001",
   "0X00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func Test_Hex(t *testing.T) {
   for _, raw_iv := range raw_ivs {
      iv, err := Segment{Raw_IV: raw_iv}.IV()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(iv)
   }
}
