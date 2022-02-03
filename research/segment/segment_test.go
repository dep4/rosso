package m3u

import (
   "fmt"
   "os"
   "testing"
)

var segments = []string{
   "segment-paramount.m3u8",
   //"segment-bbc.m3u8",
   //"segment-nbc.m3u8",
   //"segment-twitter.m3u8",
}

func TestMaster(t *testing.T) {
   for _, segment := range segments {
      fmt.Println(segment + ":")
      file, err := os.Open(segment)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      seg, err := NewSegment(file)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", seg)
   }
}
