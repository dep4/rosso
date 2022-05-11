package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

type testType struct {
   name string
   base string
}

var channel4 = testType{
   "channel4.mpd",
   "https://ak-jos-c4assets-com.akamaized.net/CH4_44_7_900_18926001001003_001/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd",
}

func TestChannel4(t *testing.T) {
   err := newMedia(channel4)
   if err != nil {
      t.Fatal(err)
   }
}

func TestRoku(t *testing.T) {
   err := newMedia(roku)
   if err != nil {
      t.Fatal(err)
   }
}

func newMedia(test testType) error {
   base, err := url.Parse(test.base)
   if err != nil {
      return err
   }
   file, err := os.Open(test.name)
   if err != nil {
      return err
   }
   defer file.Close()
   adas, err := NewAdaptationSet(file)
   if err != nil {
      return err
   }
   video := adas.MimeType(Video).Represent(0)
   init, err := video.Initialization(base)
   if err != nil {
      return err
   }
   fmt.Println(init)
   media, err := video.Media(base)
   if err != nil {
      return err
   }
   for _, addr := range media {
      fmt.Println(addr)
   }
   for _, ada := range adas.MimeType(Video) {
      for _, rep := range ada.Representation {
         fmt.Println(rep)
      }
   }
   for _, ada := range adas.MimeType(Audio) {
      for _, rep := range ada.Representation {
         fmt.Println(rep)
      }
   }
   return nil
}

