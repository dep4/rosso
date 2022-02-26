package m3u

import (
   "encoding/json"
   "fmt"
   "net/http"
   "path"
   "strings"
   "testing"
)

func (m mediaset) masters() ([]Master, error) {
   href := m.Media[1].Connection[0].Href
   fmt.Println("GET", href)
   res, err := http.Get(href)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   dir, _ := path.Split(href)
   return Decoder{dir}.Masters(res.Body)
}

func TestM3U(t *testing.T) {
   set, err := newMediaset()
   if err != nil {
      t.Fatal(err)
   }
   mass, err := set.masters()
   if err != nil {
      t.Fatal(err)
   }
   addr := mass[0].URI
   fmt.Println("GET", addr)
   res, err := http.Get(addr)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   dir, _ := path.Split(addr)
   segs := Decoder{dir}.Segments(res.Body)
   for _, seg := range segs {
      fmt.Println(seg)
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
   var buf strings.Builder
   buf.WriteString("http://open.live.bbc.co.uk")
   buf.WriteString("/mediaselector/6/select/version/2.0/mediaset/pc/vpid/")
   buf.WriteString("p0bkp6nx")
   fmt.Println("GET", buf.String())
   res, err := http.Get(buf.String())
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
