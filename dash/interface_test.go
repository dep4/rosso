package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "strings"
   "testing"
)

type test_filter struct {
   bandwidth int
}

func (test_filter) Audio(r Representations) Representations {
   return Audio(r)
}

func (test_filter) Video(r Representations) Representations {
   return Video(r)
}

func (t test_filter) Video_Index(r Representations) int {
   return Bandwidth(r, t.bandwidth)
}

func (test_filter) Audio_Index(r Representations) int {
   carry := -1
   for i, item := range r {
      if !strings.Contains(item.Codecs, "mp4a.") {
         continue
      }
      if !strings.HasPrefix(item.Adaptation.Lang, "en") {
         continue
      }
      if item.Role() == "description" {
         continue
      }
      carry = i
   }
   return carry
}

var tests = []string{
   "mpd/amc-clear.mpd",
   "mpd/amc-protected.mpd",
   "mpd/paramount-lang.mpd",
   "mpd/paramount-role.mpd",
   "mpd/roku.mpd",
}

func Test_Info(t *testing.T) {
   test := test_filter{988020}
   for _, name := range tests {
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      var pre Presentation
      if err := xml.NewDecoder(file).Decode(&pre); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      reps := pre.Representation()
      items := test.Audio(reps)
      for i, item := range items {
         if i == test.Audio_Index(items) {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      items = test.Video(reps)
      for i, item := range items {
         if i == test.Video_Index(items) {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      fmt.Println()
   }
}
