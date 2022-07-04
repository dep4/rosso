package hls

import (
   "fmt"
   "os"
   "testing"
)

type stream_filter struct {
   audio string
   codecs []string
}

func (s stream_filter) Audio() string { return s.audio }

func (s stream_filter) Codecs() []string { return s.codecs }

var stream_filters = map[string]Stream_Filter{
   "m3u8/nbc-master.m3u8": nil,
   "m3u8/roku-master.m3u8": nil,
   "m3u8/paramount-master.m3u8": &stream_filter{
      codecs: []string{"avc1."},
   },
   "m3u8/cbc-master.m3u8": &stream_filter{
      codecs: []string{"avc1."},
   },
   "m3u8/apple-master.m3u8": &stream_filter{
      audio: "-ak-",
      codecs: []string{"avc1.", "mp4a."},
   },
}

func Test_Stream(t *testing.T) {
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
