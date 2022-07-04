package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func apple(audio, codecs string) bool {
   if !strings.Contains(codecs, "avc1.") {
      return false
   }
   if !strings.Contains(codecs, "mp4a.") {
      return false
   }
   if !strings.Contains(audio, "-ak-") {
      return false
   }
   return true
}

var stream_filters = map[string]Stream_Filter{
   "m3u8/nbc-master.m3u8": nil,
   "m3u8/roku-master.m3u8": nil,
   "m3u8/paramount-master.m3u8": avc1,
   "m3u8/cbc-master.m3u8": avc1,
   "m3u8/apple-master.m3u8": apple,
}

func avc1(audio, codecs string) bool {
   return strings.Contains(codecs, "avc1.")
}

func bandwidth(a int) int {
   b := 1
   if b > a {
      return b - a
   }
   return a - b
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
      stream := master.Streams.Filter(val).Reduce(bandwidth)
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
      for _, stream := range master.Streams.Filter(val) {
         fmt.Println(stream)
      }
      fmt.Println()
   }
}
