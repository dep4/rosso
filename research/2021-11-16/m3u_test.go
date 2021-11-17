package m3u

import (
   "fmt"
   "os"
   "testing"
)

func TestM3U(t *testing.T) {
   f, err := os.Open("HLSPlaylist.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   for k, v := range newPlaylist(f) {
      fmt.Println(k, v)
   }
}
