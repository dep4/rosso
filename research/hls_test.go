package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func TestStream(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := NewScanner(file).Master(&url.URL{})
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
}
