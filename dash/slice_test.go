package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "mpd/roku.mpd",
   "mpd/paramount-role.mpd",
   "mpd/paramount-lang.mpd",
   "mpd/amc-protected.mpd",
   "mpd/amc-clear.mpd",
}

func Test_All(t *testing.T) {
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
      for _, rep := range med.Representations() {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Audio_Reduce(t *testing.T) {
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
      rep := med.Representations().Filter(Audio).Reduce(Bandwidth(0))
      fmt.Print(name, "\n", rep, "\n\n")
   }
}

func Test_Video_Filter(t *testing.T) {
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
      for _, rep := range med.Representations().Filter(Video) {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}

func Test_Video_Reduce(t *testing.T) {
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
      rep := med.Representations().Filter(Video).Reduce(Bandwidth(0))
      fmt.Print(name, "\n", rep, "\n\n")
   }
}

func Test_Audio_Filter(t *testing.T) {
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
      for _, rep := range med.Representations().Filter(Audio) {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
