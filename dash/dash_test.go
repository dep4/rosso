package dash

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "mpd/paramount-lang.mpd",
   "mpd/paramount-role.mpd",
   "mpd/roku.mpd",
}

func TestRepresent(t *testing.T) {
   for _, test := range tests {
      var med Media
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      if _, err := med.ReadFrom(file); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      for _, rep := range med.Represents(Video) {
         fmt.Println(rep)
      }
      for _, rep := range med.Represents(Audio) {
         fmt.Println(rep)
      }
      protect := med.Protection()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(protect.Default_KID)
   }
}
