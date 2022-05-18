package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func (t testType) represent() error {
   file, err := os.Open(t.name)
   if err != nil {
      return err
   }
   defer file.Close()
   adas, err := NewAdaptationSet(file)
   if err != nil {
      return err
   }
   kID, err := adas.Protection().KID()
   if err != nil {
      return err
   }
   fmt.Printf("%x\n", kID)
   iterate := func(s string) {
      for _, ada := range adas.MimeType(s) {
         for _, rep := range ada.Representation {
            fmt.Println(rep)
         }
      }
   }
   iterate(Video)
   iterate(Audio)
   return nil
}

type testType struct {
   name string
   base string
}

var channel4 = testType{
   "channel4.mpd",
   "https://ak-jos-c4assets-com.akamaized.net/CH4_44_7_900_18926001001003_001/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd",
}

var roku = testType{
   "roku.mpd",
   "https://vod.delivery.roku.com/41e834bbaecb4d27890094e3d00e8cfb/aaf72928242741a6ab8d0dfefbd662ca/87fe48887c78431d823a845b377a0c0f/index.mpd",
}

func TestChannel4Media(t *testing.T) {
   err := channel4.media()
   if err != nil {
      t.Fatal(err)
   }
}

func TestChannel4Represent(t *testing.T) {
   err := channel4.represent()
   if err != nil {
      t.Fatal(err)
   }
}

func TestRokuMedia(t *testing.T) {
   err := roku.media()
   if err != nil {
      t.Fatal(err)
   }
}

func TestRokuRepresent(t *testing.T) {
   err := roku.represent()
   if err != nil {
      t.Fatal(err)
   }
}

func (t testType) media() error {
   base, err := url.Parse(t.base)
   if err != nil {
      return err
   }
   file, err := os.Open(t.name)
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
   return nil
}
