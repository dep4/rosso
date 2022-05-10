package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

type test struct {
   name string
   base string
}

var channel4 = test{
   "channel4.mpd",
   "https://ak-jos-c4assets-com.akamaized.net/CH4_44_7_900_18926001001003_001/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd",
}

var roku = test{
   "roku.mpd",
   "https://vod.delivery.roku.com/41e834bbaecb4d27890094e3d00e8cfb/aaf72928242741a6ab8d0dfefbd662ca/87fe48887c78431d823a845b377a0c0f/index.mpd",
}

func TestMedia(t *testing.T) {
   base, err := url.Parse(roku.base)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Open(roku.name)
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   adas, err := NewAdaptationSet(file)
   if err != nil {
      t.Fatal(err)
   }
   video := adas.MimeType(Video).Represent(0)
   init, err := video.Initialization(base)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(init)
   media, err := video.Media(base)
   if err != nil {
      t.Fatal(err)
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
}
