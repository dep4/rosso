package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

/*
if !strings.HasPrefix(rep.Adaptation.Lang, "en") {
if rep.Role() == "description" {
*/
func audio(r Representation) bool {
   return r.MimeType == "audio/mp4"
}

var tests = map[string]Filter{
   "mpd/amc-clear.mpd": nil,
   "mpd/amc-protected.mpd": nil,
   "mpd/paramount-lang.mpd": nil,
   "mpd/paramount-role.mpd": nil,
   "mpd/roku.mpd": nil,
}

func Test_Audio(t *testing.T) {
   for name, callback := range tests {
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
      for _, rep := range med.Representations().Filter(callback) {
         fmt.Println(rep)
      }
      fmt.Println()
   }
}
