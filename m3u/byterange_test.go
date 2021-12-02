package m3u

import (
   "fmt"
   "os"
   "testing"
)

const byteRange = "HLS_540.m3u8"

func TestRange(t *testing.T) {
   file, err := os.Open(byteRange)
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   ran := NewByteRange(file)
   fmt.Println(ran)
}
