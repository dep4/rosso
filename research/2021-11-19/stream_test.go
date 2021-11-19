package m3u

import (
   "fmt"
   "os"
   "testing"
)

const prefix = "http://v.redd.it/16cqbkev2ci51/"

func TestStream(t *testing.T) {
   f, err := os.Open("HLS_540.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   s := newStream(f, prefix)
   fmt.Println(s)
}
