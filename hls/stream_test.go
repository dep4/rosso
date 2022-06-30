package hls

import (
   "fmt"
   "os"
   "testing"
)

func Test_Stream_Some(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   streams := master.Streams.
      Audio("-ak-").
      Audio("-stereo-").
      Codecs("dvh1")
   for _, stream := range streams {
      fmt.Println(stream)
   }
}

func Test_Stream_All(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   for _, stream := range master.Streams {
      fmt.Println(stream)
   }
}
