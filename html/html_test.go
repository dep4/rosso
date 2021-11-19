package html

import (
   "fmt"
   "os"
   "testing"
)

type release struct {
   Image string `json:"og:image"`
   Release_Date string `json:"music:release_date"`
}

func TestHTML(t *testing.T) {
   f, err := os.Open("bleep.html")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   var rel release
   if err := NewStringMap(f).Struct(&rel); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", rel)
}
