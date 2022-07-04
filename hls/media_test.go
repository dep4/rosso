package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

type media_filter struct {
   id string
   name string
   typ string
}

func (m media_filter) Group_ID(s string) bool {
   return strings.Contains(s, m.id)
}

func (m media_filter) Name(s string) bool {
   return m.name == "" || m.name == s
}

func (m media_filter) Type(s string) bool {
   return m.typ == s
}

var media_filters = map[string]Media_Filter{
   "m3u8/paramount-master.m3u8": nil, // 0
   "m3u8/roku-master.m3u8": nil, // 0
   "m3u8/nbc-master.m3u8": media_filter{typ: "AUDIO"}, // 0
   "m3u8/cbc-master.m3u8": media_filter{typ: "AUDIO"}, // 2
   "m3u8/apple-master.m3u8": media_filter{
      id: "-ak-",
      name: "English",
      typ: "AUDIO",
   },
}

func Test_Media(t *testing.T) {
   for key, val := range media_filters {
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
      for _, medium := range master.Media.Filter(val) {
         fmt.Println(medium)
      }
      fmt.Println()
   }
}
