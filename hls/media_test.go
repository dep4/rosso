package hls

import (
   "fmt"
   "os"
   "testing"
)

func Test_Media(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   for _, medium := range master.Media {
      fmt.Println(medium)
   }
}
