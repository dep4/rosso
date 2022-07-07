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
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      reps := med.Representations().Filter(Audio)
      rep := reps.Reduce(func(carry, item Representation) bool {
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
      fmt.Print(name, "\n", rep, "\n\n")
   }
}

func Test_Video(t *testing.T) {
   distance := Bandwidth(432000)
   for _, name := range tests {
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      reps := med.Representations().Filter(Video)
      rep := reps.Reduce(func(carry, item Representation) bool {
         return distance(item) < distance(carry)
      })
      fmt.Print(name, "\n", rep, "\n\n")
   }
}

func Test_Info(t *testing.T) {
   for _, name := range tests {
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for _, rep := range med.Representations().Filter(Audio_Video) {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
