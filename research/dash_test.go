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
   sets, err := AdaptationSets(file)
   if err != nil {
      t.Fatal(err)
   }
   for _, set := range sets {
      if set.Main() {
         fmt.Println(set.SegmentTemplate.Media)
         for _, rep := range set.Representation {
            fmt.Printf("%+v\n", rep)
         }
      }
   }
}
