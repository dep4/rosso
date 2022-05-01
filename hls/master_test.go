package hls

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const cbcMaster =
   "https://cbcrcott-gem.akamaized.net/0f73fb9d-87f0-4577-81d1-e6e970b89a69" +
   "/CBC_DOWNTON_ABBEY_S01E05.ism/desktop_master.m3u8"

func TestMaster(t *testing.T) {
   file, err := os.Open("m3u8/cbc-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   addr, err := url.Parse(cbcMaster)
   if err != nil {
      t.Fatal(err)
   }
   master, err := NewScanner(file).Master(addr)
   if err != nil {
      t.Fatal(err)
   }
   for i, video := range master.Stream {
      if i == 0 {
         audio := master.Audio(video)
         fmt.Printf("%+v\n", audio)
      }
      fmt.Println(video)
   }
}
