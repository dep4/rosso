package m3u

import (
   "fmt"
   "os"
   "testing"
)

const prefix = "http://v.redd.it/pu8r27nbhhl41/"

func TestPlaylist(t *testing.T) {
   f, err := os.Open("HLSPlaylist.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   for key, val := range NewPlaylist(f, prefix) {
      fmt.Print(key, "\n", val, "\n")
   }
}

func TestStream(t *testing.T) {
   f, err := os.Open("HLS_540_v4.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   s := NewStream(f, prefix)
   fmt.Println(s)
}
