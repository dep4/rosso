package m3u

import (
   "fmt"
   "os"
   "testing"
)

func TestPlaylist(t *testing.T) {
   f, err := os.Open("HLSPlaylist.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   p, err := NewPlaylist(f)
   if err != nil {
      t.Fatal(err)
   }
   for key, val := range p {
      fmt.Print(key, "\n", val, "\n")
   }
}
