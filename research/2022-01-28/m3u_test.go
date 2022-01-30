package m3u

import (
   "fmt"
   "os"
   "testing"
)

var names = []string{
   "master-bbc.m3u8",
   "master-nbc.m3u8",
   "master-paramount.m3u8",
}

func TestMaster(t *testing.T) {
   for _, name := range names {
      fmt.Println(name + ":")
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      for _, form := range decode(file) {
         fmt.Printf("%+v\n", form)
      }
   }
}
