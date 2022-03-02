package hls

import (
   "fmt"
   "net/http"
   "strings"
   "testing"
)

func platformOne() (string, error) {
   var buf strings.Builder
   buf.WriteString("http://link.theplatform.com/s/dJ5BDC/media/guid/2198311517")
   buf.WriteString("/3htV4fvVt4Z8gDZHqlzPOGLSMgcGc_vy")
   req, err := http.NewRequest("GET", buf.String(), nil)
   if err != nil {
      return "", err
   }
   req.URL.RawQuery = "formats=M3U"
   // this redirects
   res, err := new(http.Client).Do(req)
   if err != nil {
      return "", err
   }
   defer res.Body.Close()
   mas, err := one(res)
   if err != nil {
      return "", err
   }
   return mas.stream[0].URI, nil
}

func TestThree(t *testing.T) {
   href, err := platformOne()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(href)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   segs, err := two(res)
   if err != nil {
      t.Fatal(err)
   }
   for _, seg := range segs {
      fmt.Printf("%+v\n", seg)
   }
}
