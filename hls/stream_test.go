package hls

import (
   "fmt"
   "os"
   "testing"
)

func Test_Stream_Some(t *testing.T) {
   names := []string{
      "m3u8/apple-master.m3u8",
      "m3u8/cbc-master.m3u8",
      "m3u8/nbc-master.m3u8",
      "m3u8/paramount-master.m3u8",
      "m3u8/roku-master.m3u8",
   }
   for i, name := range names {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(name)
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      master, err := New_Scanner(file).Master()
      if err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      for _, stream := range master.Streams.Video() {
         fmt.Println(stream)
      }
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
