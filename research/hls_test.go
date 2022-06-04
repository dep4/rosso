package hls

import (
   "fmt"
   "github.com/89z/format/hls"
   "net/url"
   "os"
   "testing"
)

func TestMaster(t *testing.T) {
   file, err := os.Open("ignore.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := hls.NewScanner(file).Master(&url.URL{})
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
}
