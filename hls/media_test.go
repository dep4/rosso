package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func audio(m Medium) bool {
   return m.Type == "AUDIO"
}

// apple reduce
// return strings.Contains(group_ID, "-160_")

func apple_media(m Medium) bool {
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

var media_tests = map[string]Media_Filter{
   "m3u8/paramount-master.m3u8": nil,
   "m3u8/roku-master.m3u8": nil,
   "m3u8/nbc-master.m3u8": audio,
   "m3u8/cbc-master.m3u8": audio,
   "m3u8/apple-master.m3u8": apple_media,
}

func Test_Media_Filter(t *testing.T) {
   for name, callback := range media_tests {
      fmt.Println(name)
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
      for _, medium := range master.Media.Filter(callback) {
         fmt.Println(medium)
      }
      fmt.Println()
   }
}

func Test_Media_All(t *testing.T) {
   for name := range media_tests {
      fmt.Println(name)
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
      for _, medium := range master.Media {
         fmt.Println(medium)
      }
      fmt.Println()
   }
}
