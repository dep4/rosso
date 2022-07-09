package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "strings"
   "testing"
)

var tests = []string{
   "mpd/amc-clear.mpd",
   "mpd/amc-protected.mpd",
   "mpd/paramount-lang.mpd",
   "mpd/paramount-role.mpd",
   "mpd/roku.mpd",
}

func Test_Audio(t *testing.T) {
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
      reps := pre.Representation().Filter(Audio)
      target := reps.Index(func(carry, item Representation) bool {
         if !strings.HasPrefix(item.Adaptation.Lang, "en") {
            return false
         }
         if !strings.Contains(item.Codecs, "mp4a.") {
            return false
         }
         if item.Role() == "description" {
            return false
         }
         return true
      })
      fmt.Println(name)
      for i, rep := range reps {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Video(t *testing.T) {
   distance := Bandwidth(0)
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
      reps := pre.Representation().Filter(Video)
      target := reps.Index(func(carry, item Representation) bool {
         return distance(item) < distance(carry)
      })
      fmt.Println(name)
      for i, rep := range reps {
         if i == target {
            fmt.Print("!")
         }
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Info(t *testing.T) {
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
      reps := pre.Representation().Filter(func(r Representation) bool {
         return Audio(r) || Video(r)
      })
      for _, rep := range reps {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
