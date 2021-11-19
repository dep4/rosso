package m3u

import (
   "fmt"
   "os"
   "testing"
)

const prefix = "http://v.redd.it/16cqbkev2ci51/"

func TestPlaylist(t *testing.T) {
   f, err := os.Open("HLSPlaylist.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   p, err := NewPlaylist(f, prefix)
   if err != nil {
      t.Fatal(err)
   }
   for key, val := range p {
      fmt.Print(key, "\n", val, "\n")
   }
}

func TestStream(t *testing.T) {
   f, err := os.Open("HLS_540.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   s := NewStream(f, prefix)
   fmt.Println(s)
}
