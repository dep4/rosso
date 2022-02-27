package m3u

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "cbs.m3u8",
   "nbc.m3u8",
}

func TestPlaylist(t *testing.T) {
   for _, test := range tests {
      fmt.Println(test + ":")
      buf, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      for _, form := range Unmarshal(buf, "http://example.com/") {
         fmt.Printf("%+v\n", form)
      }
      fmt.Println()
   }
}
