package m3u

import (
   "fmt"
   "os"
   "testing"
)

var segments = []string{
   "segment-bbc.m3u8",
   "segment-nbc.m3u8",
   "segment-paramount.m3u8",
   "segment-twitter.m3u8",
}

func TestSegment(t *testing.T) {
   for _, segment := range segments {
      file, err := os.Open(segment)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      seg, err := NewSegment(file)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(segment + ":")
      fmt.Printf("%q\n", seg.Key)
      for _, addr := range seg.URI {
         fmt.Println(addr)
      }
   }
}
