package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func (s stream_filter) Audio(v string) bool {
   return strings.Contains(v, s.audio)
}

func (s stream_filter) Bandwidth(v int) int {
   if v > s.bandwidth {
      return v - s.bandwidth
   }
   return s.bandwidth - v
}

type stream_filter struct {
   audio string
   bandwidth int
   codecs []string
}

func (s stream_filter) Codecs(v string) bool {
   for _, curr := range s.codecs {
      if !strings.Contains(v, curr) {
         return false
      }
   }
   return true
}

var stream_filters = map[string]Stream_Filter{
   "m3u8/nbc-master.m3u8": stream_filter{},
   "m3u8/roku-master.m3u8": stream_filter{},
   "m3u8/paramount-master.m3u8": stream_filter{
      codecs: []string{"avc1."},
   },
   "m3u8/cbc-master.m3u8": stream_filter{
      codecs: []string{"avc1."},
   },
   "m3u8/apple-master.m3u8": stream_filter{
      audio: "-ak-",
      codecs: []string{"avc1.", "mp4a."},
   },
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
      stream := master.Streams.Streams(val).Bandwidth(val)
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
