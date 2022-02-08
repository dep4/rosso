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
