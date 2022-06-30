package hls

import (
   "fmt"
   "os"
   "testing"
)

func Test_Media_Some(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   media := master.Media.
      Group_ID("-ak-").
      Group_ID("-stereo-").
      Name("English").
      Type("AUDIO")
   fmt.Println(media)
}

func Test_Media_All(t *testing.T) {
   file, err := os.Open("ignore/apple-master.m3u8")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   master, err := New_Scanner(file).Master()
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(master.Media)
}
