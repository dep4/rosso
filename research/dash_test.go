package dash

import (
   "fmt"
   "net/url"
   "os"
   "testing"
)

const mpd = "https://ak-jos-c4assets-com.akamaized.net" +
   "/CH4_44_7_900_18926001001003_001" +
   "/CH4_44_7_900_18926001001003_001_J01.ism/stream.mpd"

func TestDASH(t *testing.T) {
   addr, err := url.Parse(mpd)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Open("18926-001.mpd")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   adas, err := Adaptations(addr, file)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(adas)
}
