package m3u

import (
   "fmt"
   "os"
   "testing"
)

type masterTest struct {
   base, dir string
}

var masters = []masterTest{
   {base: "master-bbc.m3u8", dir: "http://example.com/"},
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
      forms, err := Decode(file, master.dir)
      if err != nil {
         t.Fatal(err)
      }
      for _, form := range forms {
         fmt.Printf("%q %q %.9q\n", form.resolution, form.codecs, form.uri)
      }
   }
}
