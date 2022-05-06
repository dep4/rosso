package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

func TestRepresent(t *testing.T) {
   body, err := os.Open("roku.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer body.Close()
   period, err := NewPeriod(body)
   if err != nil {
      t.Fatal(err)
   }
   video := period.Video(2035878)
   fmt.Printf("%a\n", video)
   audio := period.Audio(128000)
   fmt.Printf("%a\n", audio)
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
            addrs, err := ada.SegmentTemplate.URL(rep, base)
            if err != nil {
               t.Fatal(err)
            }
            for _, addr := range addrs {
               fmt.Println(addr)
            }
            fmt.Println(ada.SegmentTemplate.Base(rep))
         }
      }
   }
}
