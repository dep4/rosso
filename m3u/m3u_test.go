package m3u

import (
   "fmt"
   "os"
   "testing"
)

const (
   byteRange = "iduchl/HLS_540.m3u8"
   prefix = "http://v.redd.it/pu8r27nbhhl41/"
   stream = "fffrnw/HLSPlaylist.m3u8"
)

func TestRange(t *testing.T) {
   f, err := os.Open(byteRange)
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   b := NewByteRange(f, prefix)
   fmt.Println(b)
}

func TestStream(t *testing.T) {
   f, err := os.Open(stream)
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   for _, dir := range Streams(f, prefix) {
      fmt.Println(dir)
   }
}
