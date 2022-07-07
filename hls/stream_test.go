package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

var tests = []string{
   "m3u8/nbc-master.m3u8",
   "m3u8/roku-master.m3u8",
   "m3u8/paramount-master.m3u8",
   "m3u8/cbc-master.m3u8",
   "m3u8/apple-master.m3u8",
}

func Test_Stream(t *testing.T) {
   for _, name := range tests {
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
      fmt.Println(name)
      for _, stream := range master.Stream {
         fmt.Println(stream)
      }
      fmt.Println()
   }
}
