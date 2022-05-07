package dash

import (
   "fmt"
   "net/http"
   "testing"
)

var tests = map[string]string{
   //"channel4": "https://ak-jos-c4assets-com.akamaized.net/CH4_44_7_900_18926001001003_001/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd",
   "roku": "https://vod.delivery.roku.com/41e834bbaecb4d27890094e3d00e8cfb/aaf72928242741a6ab8d0dfefbd662ca/87fe48887c78431d823a845b377a0c0f/index.mpd",
}

func TestDASH(t *testing.T) {
   for _, test := range tests {
      fmt.Println("GET", test)
      res, err := http.Get(test)
      if err != nil {
         t.Fatal(err)
      }
      period, err := NewPeriod(res.Body)
      if err != nil {
         t.Fatal(err)
      }
      video := period.Video(0)
      addrs, err := video.URL(res.Request.URL)
      if err != nil {
         t.Fatal(err)
      }
      for _, addr := range addrs {
         fmt.Println(addr)
      }
      fmt.Println(video.Base())
      if err := res.Body.Close(); err != nil {
         t.Fatal(err)
      }
   }
}
