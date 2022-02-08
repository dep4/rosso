package m3u

import (
   "fmt"
   "os"
   "testing"
)

type testType struct {
   base, dir string
}

const dir = "http://example.com/"

var masters = []testType{
   {base: "master-bbc.m3u8", dir: dir},
   {base: "master-nbc.m3u8"},
   {base: "master-paramount.m3u8"},
}

var segments = []testType{
   {base: "segment-bbc.m3u8", dir: dir},
   {base: "segment-nbc.m3u8"},
   {base: "segment-paramount.m3u8"},
   {base: "segment-twitter.m3u8", dir: dir},
}

func TestMaster(t *testing.T) {
   for _, master := range masters {
      fmt.Println(master.base + ":")
      file, err := os.Open(master.base)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      mass, err := Decoder{master.dir}.Masters(file)
      if err != nil {
         t.Fatal(err)
      }
      for _, mas := range mass {
         fmt.Printf("%+v\n", mas)
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
      seg, err := Decoder{segment.dir}.Segment(file)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", seg)
   }
}
