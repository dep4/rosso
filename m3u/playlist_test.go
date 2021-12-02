package m3u

import (
   "fmt"
   "os"
   "testing"
)

const playlist = "pc_hd_abr_v2_hls_master.m3u8"

func TestPlaylist(t *testing.T) {
   file, err := os.Open(playlist)
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   for key, val := range NewPlaylist(file) {
      fmt.Print(key, "\n", val, "\n")
   }
}
