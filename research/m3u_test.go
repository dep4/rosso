package m3u

import (
   "encoding/json"
   "fmt"
   "net/http"
   "strings"
   "testing"
)

func TestM3U(t *testing.T) {
   set, err := newMediaset()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(set.Media[1].Connection[0].Href)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   mass, err := Masters(res.Body)
   if err != nil {
      t.Fatal(err)
   }
   for _, mas := range mass {
      fmt.Printf("%+v\n", mas)
   }
}

type mediaset struct {
   Media []struct {
      Connection []struct {
         Href string
      }
   }
}

func newMediaset() (*mediaset, error) {
   var str strings.Builder
   str.WriteString("http://open.live.bbc.co.uk")
   str.WriteString("/mediaselector/6/select/version/2.0/mediaset/pc/vpid/")
   str.WriteString("p0bkp6nx")
   res, err := http.Get(str.String())
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   set := new(mediaset)
   if err := json.NewDecoder(res.Body).Decode(set); err != nil {
      return nil, err
   }
   return set, nil
}
