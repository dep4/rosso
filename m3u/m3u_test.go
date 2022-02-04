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

func TestMaster(t *testing.T) {
   for _, master := range masters {
      fmt.Println(master.base + ":")
      file, err := os.Open(master.base)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      var buf Scanner
      buf.Init(file)
      for buf.Scan() {
         fmt.Printf("%+v\n", buf.Master)
      }
   }
}
