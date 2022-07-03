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
      seg := New_Scanner(file).Segment()
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n\n", seg)
   }
}
