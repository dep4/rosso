package dash

import (
   "encoding/hex"
   "fmt"
   "net/http"
   "net/url"
   "os"
   "testing"
)

const hexKey = "13d7c7cf295444944b627ef0ad2c1b3c"

func newAdaptationSet() (*url.URL, AdaptationSet, error) {
   res, err := http.Get(roku.base)
   if err != nil {
      return nil, nil, err
   }
   defer res.Body.Close()
   adas, err := NewAdaptationSet(res.Body)
   if err != nil {
      return nil, nil, err
   }
   return res.Request.URL, adas, nil
}

func TestDecrypt(t *testing.T) {
   dst, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dst.Close()
   base, adas, err := newAdaptationSet()
   if err != nil {
      t.Fatal(err)
   }
   audio := adas.MimeType(Audio).Represent(0)
   init, err := audio.Initialization(base)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(init)
   res, err := http.Get(init.String())
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   dst.ReadFrom(res.Body)
   media, err := audio.Media(base)
   if err != nil {
      t.Fatal(err)
   }
   key, err := hex.DecodeString(hexKey)
   if err != nil {
      t.Fatal(err)
   }
   for _, addr := range media[:9] {
      fmt.Println(addr)
      res, err := http.Get(addr.String())
      if err != nil {
         t.Fatal(err)
      }
      if res.StatusCode != http.StatusOK {
         t.Fatal(res.Status, addr)
      }
      if err := Decrypt(dst, res.Body, key); err != nil {
         t.Fatal(err)
      }
      res.Body.Close()
   }
}

var roku = testType{
   "roku.mpd",
   "https://vod.delivery.roku.com/41e834bbaecb4d27890094e3d00e8cfb/aaf72928242741a6ab8d0dfefbd662ca/87fe48887c78431d823a845b377a0c0f/index.mpd",
}
