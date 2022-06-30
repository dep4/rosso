package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "mpd/paramount-lang.mpd",
   "mpd/paramount-role.mpd",
   "mpd/roku.mpd",
}

func Test_Filter(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
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
      reps := med.Representations()
      fmt.Println(test, "video")
      for _, rep := range reps.Filter_Codecs("avc1") {
         fmt.Println(rep)
      }
      fmt.Println(test, "audio")
      rep := reps.Reduce_Codecs("mp4a")
      fmt.Println(rep) 
   }
}

func Test_Representations(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
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
      fmt.Println(test)
      for _, rep := range med.Representations() {
         fmt.Println(rep)
      }
   }
}
