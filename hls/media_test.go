package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func audio_media_filter(group_ID, name, typ string) bool {
   return typ == "AUDIO"
}

func apple_media_filter(group_ID, name, typ string) bool {
   if !strings.Contains(group_ID, "-ak-") {
      return false
   }
   if name != "English" {
      return false
   }
   if typ != "AUDIO" {
      return false
   }
   return true
}

type filter_reduce struct {
   filter func(string, string, string) bool
   reduce func(string, string) bool
}

func cbc_media_reduce(group_ID, name string) bool {
   return name == "English"
}

func apple_media_reduce(group_ID, name string) bool {
   return strings.Contains(group_ID, "-160_")
}

var media_tests = map[string]filter_reduce {
   "m3u8/paramount-master.m3u8": {},
   "m3u8/roku-master.m3u8": {},
   "m3u8/nbc-master.m3u8": {
      filter: audio_media_filter,
   },
   "m3u8/cbc-master.m3u8": {
      filter: audio_media_filter,
      reduce: cbc_media_reduce,
   },
   "m3u8/apple-master.m3u8": {
      filter: apple_media_filter,
      reduce: apple_media_reduce,
   },
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
      medium := master.Media.Media(val.filter).Medium(val.reduce)
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
      for _, medium := range master.Media.Filter(val.filter) {
         fmt.Println(medium)
      }
      fmt.Println()
   }
}
