package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func bandwidth(v int) Reduce {
   distance := func(r *Representation) int {
      if r.Bandwidth > v {
         return r.Bandwidth - v
      }
      return v - r.Bandwidth
   }
   return func(carry *Representation, item Representation) *Representation {
      if carry == nil || distance(&item) < distance(carry) {
         return &item
      }
      return carry
   }
}

var tests = []string{
   "mpd/roku.mpd",
   "mpd/paramount-role.mpd",
   "mpd/paramount-lang.mpd",
   "mpd/amc-protected.mpd",
   "mpd/amc-clear.mpd",
}

func Test_Video(t *testing.T) {
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
      rep := med.Representations().Filter(Video).Reduce(bandwidth(0))
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
      for _, rep := range med.Representations().Filter(AudioVideo) {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
