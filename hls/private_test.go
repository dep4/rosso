package hls

import (
   "fmt"
   "testing"
)

var ivs = []string{
   "00000000000000000000000000000001",
   "0X00000000000000000000000000000001",
   "0x00000000000000000000000000000001",
}

func TestHex(t *testing.T) {
   for _, iv := range ivs {
      buf, err := scanHex(iv)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(buf)
   }
}
