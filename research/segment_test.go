package m3u

import (
   "fmt"
   "os"
   "testing"
)

var segments = []testType{
   {base: "segment-bbc.m3u8", dir: dir},
   {base: "segment-nbc.m3u8"},
   {base: "segment-paramount.m3u8"},
   {base: "segment-twitter.m3u8", dir: dir},
}

func TestSegment(t *testing.T) {
   for _, segment := range segments {
      file, err := os.Open(segment.base)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      seg, err := NewSegment(file)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(segment.base + ":")
      fmt.Printf("%.99q\n", seg.Key)
      for _, addr := range seg.URI {
         fmt.Printf("%.99q\n", addr)
      }
   }
}
