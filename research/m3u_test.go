package m3u

import (
   "fmt"
   "os"
   "testing"
)

var masters = []string{
   "master-bbc.m3u8",
   "master-nbc.m3u8",
   "master-paramount.m3u8",
}

func TestMaster(t *testing.T) {
   for _, master := range masters {
      fmt.Println(master + ":")
      file, err := os.Open(master)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      buf := Scanner{Reader: file}
      for buf.Scan() {
         fmt.Println(buf.Master)
      }
   }
}
