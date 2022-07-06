package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func avc1_stream_filter(s Stream) bool {
   return strings.Contains(s.Codecs, "avc1.")
}

func apple_stream_filter(s Stream) bool {
   if !strings.Contains(s.Audio, "-ak-") {
      return false
   }
   if !strings.Contains(s.Codecs, "avc1.") {
      return false
   }
   if !strings.Contains(s.Codecs, "mp4a.") {
      return false
   }
   return true
}

var stream_filters = map[string]Filter[Stream]{
   "m3u8/nbc-master.m3u8": nil,
   "m3u8/roku-master.m3u8": nil,
   "m3u8/paramount-master.m3u8": avc1_stream_filter,
   "m3u8/cbc-master.m3u8": avc1_stream_filter,
   "m3u8/apple-master.m3u8": apple_stream_filter,
}

func Test_Stream(t *testing.T) {
   for name, callback := range stream_filters {
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
      stream := master.Streams.Filter(callback).Reduce(Bandwidth(0))
      fmt.Print(name, "\n", stream, "\n\n")
   }
}
