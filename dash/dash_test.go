package dash

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "paramount-lang.mpd",
   "paramount-role.mpd",
   "roku.mpd",
}

func TestRepresent(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      period, err := NewPeriod(file)
      if err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      kID, err := period.Protection().KID()
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%x\n", kID)
      for _, rep := range period.Represents(Video) {
         fmt.Println(rep)
      }
      for _, rep := range period.Represents(Audio) {
         fmt.Println(rep)
      }
   }
}
