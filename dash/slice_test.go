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
