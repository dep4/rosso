package m3u

import (
   "fmt"
   "os"
   "testing"
)

const (
   byteRange = "iduchl/HLS_540.m3u8"
   playlist = "fffrnw/HLSPlaylist.m3u8"
)

func TestPlaylist(t *testing.T) {
   f, err := os.Open(playlist)
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   for key, val := range NewPlaylist(f) {
      fmt.Print(key, "\n", val, "\n")
   }
}

func TestRange(t *testing.T) {
   f, err := os.Open(byteRange)
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   b := NewByteRange(f)
   fmt.Println(b)
}
