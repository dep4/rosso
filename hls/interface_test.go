package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func (h hello) Audio_Index(m []Media) int {
   carry := -1
   for i, item := range m {
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
      carry = i
   }
   return carry
}

func (hello) Video_Index(s []Stream) int {
   return Bandwidth(s, 1_599_999)
}

var tests = []hello{
   {"m3u8/nbc-master.m3u8"},
   {"m3u8/roku-master.m3u8"},
   {"m3u8/paramount-master.m3u8"},
   {"m3u8/cbc-master.m3u8"},
   {"m3u8/apple-master.m3u8"},
}

func (h hello) Video(s []Stream) []Stream {
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
   return carry
}

func (h hello) Audio(m []Media) []Media {
   var carry []Media
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
      carry = append(carry, item)
   }
   return carry
}

type hello struct {
   name string
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
      audio := test.Audio(master.Media)
      for i, item := range audio {
         if i == test.Audio_Index(audio) {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      video := test.Video(master.Stream)
      for i, item := range video {
         if i == test.Video_Index(video) {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      fmt.Println()
   }
}
