package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func audio_media_filter(m Medium) bool {
   return m.Type == "AUDIO"
}

func apple_media_filter(m Medium) bool {
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

type filter_reduce struct {
   Media_Filter
   Media_Reduce
}

func apple_media_reduce(carry *Medium, item Medium) *Medium {
   if strings.Contains(item.Group_ID, "-160_") {
      return &item
   }
   return carry
}

func cbc_media_reduce(carry *Medium, item Medium) *Medium {
   if item.Name == "English" {
      return &item
   }
   return carry
}

var media_tests = map[string]filter_reduce{
   "m3u8/paramount-master.m3u8": {nil, nil},
   "m3u8/roku-master.m3u8": {nil, nil},
   "m3u8/nbc-master.m3u8": {audio_media_filter, nil},
   "m3u8/cbc-master.m3u8": {audio_media_filter, cbc_media_reduce},
   "m3u8/apple-master.m3u8": {apple_media_filter, apple_media_reduce},
}

func Test_Media_Reduce(t *testing.T) {
   for key, val := range media_tests {
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
      medium := master.Media.Filter(val.Media_Filter).Reduce(val.Media_Reduce)
      fmt.Print(key, "\n", medium, "\n\n")
   }
}

func Test_Media_Filter(t *testing.T) {
   for key, val := range media_tests {
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
      for _, medium := range master.Media.Filter(val.Media_Filter) {
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
