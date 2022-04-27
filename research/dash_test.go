package dash

import (
   "fmt"
   "os"
   "testing"
)

func TestDASH(t *testing.T) {
   file, err := os.Open("18926-001.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   media, err := NewMedia(file)
   if err != nil {
      t.Fatal(err)
   }
   for _, ada := range media.Main() {
      for _, rep := range ada.Representation {
         fmt.Printf("%+v\n", rep)
      }
   }
}
