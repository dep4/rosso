package hls

import (
   "fmt"
   "os"
   "testing"
)

func Test_Stream_Some(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   streams := master.Streams.
      Audio("-ak-").
      Audio("-stereo-").
      Codecs("dvh1")
   for _, stream := range streams {
      fmt.Println(stream)
   }
}

func Test_Stream_All(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
}

func Test_Segment(t *testing.T) {
   file, err := os.Open("ignore/apple-audio.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   seg, err := New_Scanner(file).Segment()
   if err != nil {
      t.Fatal(err)
   }
   for _, pro := range seg.Protected {
      fmt.Println(pro)
   }
   fmt.Println(seg.Key)
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
