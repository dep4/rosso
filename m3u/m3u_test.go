package m3u

import (
   "fmt"
   "os"
   "testing"
)

const prefix = "http://v.redd.it/16cqbkev2ci51/"

var tests = []string{"HLSPlaylist.m3u8", "HLS_540.m3u8"}

func TestM3U(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      for key, val := range Parse(file, prefix) {
         fmt.Println(key)
         fmt.Println(val)
      }
   }
}
