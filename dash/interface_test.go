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

func (t test_filter) Video(r []Representation) ([]Representation, int) {
   r = Video(r)
   return r, Bandwidth(r, t.bandwidth)
}

func (test_filter) Audio(r []Representation) ([]Representation, int) {
   r = Audio(r)
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
   return r, carry
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
      items, i := test.Audio(reps)
      for j, item := range items {
         if j == i {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      items, i = test.Video(reps)
      for j, item := range items {
         if j == i {
            fmt.Print("!")
         }
         fmt.Println(item)
      }
      fmt.Println()
   }
}
