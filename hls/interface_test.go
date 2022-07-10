package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

var tests = []hello{
   {"m3u8/nbc-master.m3u8"},
   {"m3u8/roku-master.m3u8"},
   {"m3u8/paramount-master.m3u8"},
   {"m3u8/cbc-master.m3u8"},
   {"m3u8/apple-master.m3u8"},
}

type hello struct {
   name string
}

func (h hello) Video(s []Stream) ([]Stream, int) {
   var carry []Stream
   for _, item := range s {
      switch h.name {
      case "m3u8/apple-master.m3u8":
         if !strings.Contains(item.Audio, "-ak-") {
            continue
         }
         if !strings.Contains(item.Codecs, "avc1.") {
            continue
         }
         if !strings.Contains(item.Codecs, "mp4a.") {
            continue
         }
      case "m3u8/cbc-master.m3u8":
         if item.Resolution == "" {
            continue
         }
      case "m3u8/paramount-master.m3u8":
         if item.Resolution == "" {
            continue
         }
      }
      carry = append(carry, item)
   }
   return carry, Bandwidth(carry, 1_599_999)
}

func (h hello) Audio(m []Media) ([]Media, int) {
   var items []Media
   for _, item := range m {
      if item.Type != "AUDIO" {
         continue
      }
      if strings.Contains(h.name, "apple") {
         if !strings.Contains(item.Group_ID, "-ak-") {
            continue
         }
         if item.Name != "English" {
            continue
         }
      }
      items = append(items, item)
   }
   i := -1
   for j, item := range items {
      switch h.name {
      case "m3u8/apple-master.m3u8":
         if !strings.Contains(item.Group_ID, "-160_") {
            continue
         }
      case "m3u8/cbc-master.m3u8":
         if item.Name != "English" {
            continue
         }
      }
      i = j
   }
   return items, i
}

func Test_Info(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test.name)
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
      fmt.Println(test.name)
      audio, i := test.Audio(master.Media)
      for j, item := range audio {
         if j == i {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      video, i := test.Video(master.Stream)
      for j, item := range video {
         if j == i {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      fmt.Println()
   }
}
