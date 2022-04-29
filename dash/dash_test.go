package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const stream = "https://ak-jos-c4assets-com.akamaized.net" +
   "/CH4_44_7_900_18926001001003_001" +
   "/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd"

func TestDASH(t *testing.T) {
   addr, err := url.Parse(stream)
   if err != nil {
      t.Fatal(err)
   }
   src, err := os.Open("18926-001.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer src.Close()
   adas, err := Adaptations(src)
   if err != nil {
      t.Fatal(err)
   }
   for _, ada := range adas {
      fmt.Println(ada.MimeType)
      for _, rep := range ada.Representation {
         if rep.ID == "video=501712" {
            addrs, err := ada.URLs(addr, rep)
            if err != nil {
               t.Fatal(err)
            }
            for _, addr := range addrs {
               fmt.Println(addr)
            }
         }
      }
   }
}
