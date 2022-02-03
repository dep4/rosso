package m3u

import (
   "fmt"
   "os"
   "testing"
)

type testType struct {
   base, dir string
}

var masters = []testType{
   {base: "master-bbc.m3u8", dir: "http://example.com/"},
   {base: "master-nbc.m3u8"},
   {base: "master-paramount.m3u8"},
}

var segments = []testType{
   {base: "segment-bbc.m3u8"},
   {base: "segment-nbc.m3u8"},
   {base: "segment-paramount.m3u8"},
   {base: "segment-twitter.m3u8"},
}

func TestMaster(t *testing.T) {
   for _, master := range masters {
      fmt.Println(master.base + ":")
      file, err := os.Open(master.base)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      mass, err := Masters(file)
      if err != nil {
         t.Fatal(err)
      }
      for _, mas := range mass {
         fmt.Printf("%q %q %.9q\n", mas.Resolution, mas.Codecs, mas.URI)
      }
   }
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
      fmt.Printf("%q\n", seg.Key)
      for _, addr := range seg.URI {
         fmt.Println(addr)
      }
   }
}
