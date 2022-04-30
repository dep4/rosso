package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func TestRepresentation(t *testing.T) {
   body, err := os.Open("18926-001.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer body.Close()
   period, err := NewPeriod(body)
   if err != nil {
      t.Fatal(err)
   }
   video := period.Video()
   for _, rep := range video.Representation {
      fmt.Println(rep)
   }
   fmt.Printf("%+v\n", period.Audio(video))
}

const stream = "https://ak-jos-c4assets-com.akamaized.net" +
   "/CH4_44_7_900_18926001001003_001" +
   "/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd"

func TestTemplate(t *testing.T) {
   base, err := url.Parse(stream)
   if err != nil {
      t.Fatal(err)
   }
   body, err := os.Open("18926-001.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer body.Close()
   per, err := NewPeriod(body)
   if err != nil {
      t.Fatal(err)
   }
   for _, ada := range per.AdaptationSet {
      for _, rep := range ada.Representation {
         if rep.ID == "video=501712" {
            temp := ada.SegmentTemplate.Replace(rep)
            addrs, err := temp.URLs(base)
            if err != nil {
               t.Fatal(err)
            }
            for _, addr := range addrs {
               fmt.Println(addr)
            }
            fmt.Println(temp.Initialization)
         }
      }
   }
}
