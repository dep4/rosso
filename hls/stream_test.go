package hls

import (
   "fmt"
   "os"
   "testing"
)

var master_names = map[string]Reducer{
   "m3u8/nbc-master.m3u8": nil,
   "m3u8/paramount-master.m3u8": AVC1{},
   "m3u8/roku-master.m3u8": nil,
   "m3u8/cbc-master.m3u8": AVC1{},
   "m3u8/apple-master.m3u8": AVC1{},
}

func Test_Stream_Some(t *testing.T) {
   for name, red := range master_names {
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
      for _, stream := range master.Streams.Reduce(red) {
         fmt.Println(stream)
      }
      fmt.Println()
   }
}

func Test_Stream_All(t *testing.T) {
   for name := range master_names {
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
      for _, stream := range master.Streams {
         fmt.Println(stream)
      }
      fmt.Println()
   }
}
