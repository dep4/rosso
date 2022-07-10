package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

var tests = map[string]filters{
   "m3u8/nbc-master.m3u8": {nbc_media, nil},
   "m3u8/roku-master.m3u8": {nil, nil},
   "m3u8/paramount-master.m3u8": {nil, paramount_stream},
   "m3u8/cbc-master.m3u8": {cbc_media, cbc_stream},
   "m3u8/apple-master.m3u8": {apple_media, apple_stream},
}

func Test_Info(t *testing.T) {
   for key, val := range tests {
      file, err := os.Open(key)
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
      fmt.Println(key)
      for _, item := range master.Streams.Filter(val.stream) {
         fmt.Println(item)
      }
      for _, item := range master.Media.Filter(val.medium) {
         fmt.Println(item)
      }
      fmt.Println()
   }
}
