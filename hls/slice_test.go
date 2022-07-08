package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func apple_media(m Media) bool {
   if !strings.Contains(m.Group_ID, "-ak-") {
      return false
   }
   if m.Name != "English" {
      return false
   }
   if m.Type != "AUDIO" {
      return false
   }
   return true
}

func apple_stream(s Stream) bool {
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

func cbc_media(m Media) bool {
   return m.Type == "AUDIO"
}

func cbc_stream(s Stream) bool {
   return strings.Contains(s.Codecs, "avc1.")
}

func nbc_media(m Media) bool {
   return m.Type == "AUDIO"
}

func paramount_stream(s Stream) bool {
   return strings.Contains(s.Codecs, "avc1.")
}

var tests = map[string]filters{
   "m3u8/nbc-master.m3u8": {nbc_media, nil},
   "m3u8/roku-master.m3u8": {nil, nil},
   "m3u8/paramount-master.m3u8": {nil, paramount_stream},
   "m3u8/cbc-master.m3u8": {cbc_media, cbc_stream},
   "m3u8/apple-master.m3u8": {apple_media, apple_stream},
}

type filters struct {
   media func(Media) bool
   stream func(Stream) bool
}

func Test_Media(t *testing.T) {
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
      master.Media = master.Media.Filter(val.media)
      target := master.Media.Index(func(carry, item Media) bool {
         return item.Name == "English"
      })
      fmt.Println(key)
      for i, media := range master.Media {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(media)
      }
      fmt.Println()
   }
}

func Test_Stream(t *testing.T) {
   distance := Bandwidth(0)
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
      master.Stream = master.Stream.Filter(val.stream)
      target := master.Stream.Index(func(carry, item Stream) bool {
         return distance(item) < distance(carry)
      })
      fmt.Println(key)
      for i, stream := range master.Stream {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(stream)
      }
      fmt.Println()
   }
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
      for _, item := range master.Stream.Filter(val.stream) {
         fmt.Println(item)
      }
      for _, item := range master.Media.Filter(val.media) {
         fmt.Println(item)
      }
      fmt.Println()
   }
}
