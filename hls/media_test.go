package hls

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

type media_filter struct {
   group_ID string
   name string
   typ string
}

func (m media_filter) Group_ID(v string) bool {
   return m.group_ID == "" || strings.Contains(v, m.group_ID)
}

func (m media_filter) Name(v string) bool {
   return m.name == "" || m.name == v
}

func (m media_filter) Type(v string) bool {
   return m.typ == "" || m.typ == v
}

var media_tests = map[string]Media_Filter{
   "m3u8/paramount-master.m3u8": media_filter{},
   "m3u8/roku-master.m3u8": media_filter{},
   "m3u8/nbc-master.m3u8": media_filter{typ: "AUDIO"},
   "m3u8/cbc-master.m3u8": media_filter{
      name: "English",
      typ: "AUDIO",
   },
   "m3u8/apple-master.m3u8": media_filter{
      // return strings.Contains(group_ID, "-160_")
      group_ID: "-ak-",
      name: "English",
      typ: "AUDIO",
   },
}

func Test_Media_Filter(t *testing.T) {
   for key := range media_tests {
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
      for _, medium := range master.Media {
         fmt.Println(medium)
      }
      fmt.Println()
   }
}
