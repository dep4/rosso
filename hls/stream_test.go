package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func apple_stream_filter(bandwidth int, audio, codecs string) int {
   if !strings.Contains(codecs, "avc1.") {
      return 0
   }
   if !strings.Contains(codecs, "mp4a.") {
      return 0
   }
   if !strings.Contains(audio, "-ak-") {
      return 0
   }
   return 1
}

var stream_filters = map[string]Stream_Func{
   "m3u8/nbc-master.m3u8": nil,
   "m3u8/roku-master.m3u8": nil,
   "m3u8/paramount-master.m3u8": avc1_stream_filter,
   "m3u8/cbc-master.m3u8": avc1_stream_filter,
   "m3u8/apple-master.m3u8": apple_stream_filter,
}

func avc1_stream_filter(bandwidth int, audio, codecs string) int {
   if strings.Contains(codecs, "avc1.") {
      return 1
   }
   return 0
}

func bandwidth_stream_reduce(bandwidth int, audio, codecs string) int {
   b := 1
   if bandwidth > b {
      return bandwidth - b
   }
   return b - bandwidth
}

func Test_Stream_Reduce(t *testing.T) {
   for key, val := range stream_filters {
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
      stream := master.Streams.Streams(val).Stream(bandwidth_stream_reduce)
      fmt.Print(key, "\n", stream, "\n\n")
   }
}

func Test_Stream_Filter(t *testing.T) {
   for key, val := range stream_filters {
      fmt.Println(key)
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
      for _, stream := range master.Streams.Streams(val) {
         fmt.Println(stream)
      }
      fmt.Println()
   }
}
