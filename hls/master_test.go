package hls

import (
   "fmt"
   "os"
   "testing"
)

func TestStreams(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := NewScanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
}

func TestMedia(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := NewScanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   for _, medium := range master.Media {
      fmt.Println(medium)
   }
}
