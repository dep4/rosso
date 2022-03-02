package hls

import (
   "fmt"
   "net/http"
   "strings"
   "testing"
)

func platform() (string, error) {
   var buf strings.Builder
   buf.WriteString("http://link.theplatform.com/s/dJ5BDC/media/guid/2198311517")
   buf.WriteString("/3htV4fvVt4Z8gDZHqlzPOGLSMgcGc_vy")
   req, err := http.NewRequest("GET", buf.String(), nil)
   if err != nil {
      return "", err
   }
   // We need "MPEG4", otherwise you get a "EXT-X-KEY" with "skd" scheme:
   req.URL.RawQuery = "formats=MPEG4,M3U"
   // This redirects:
   res, err := new(http.Client).Do(req)
   if err != nil {
      return "", err
   }
   defer res.Body.Close()
   mas, err := newMaster(res)
   if err != nil {
      return "", err
   }
   return mas.stream[0].URI, nil
}

func TestSegment(t *testing.T) {
   href, err := platform()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(href)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   seg, err := newSegment(res)
   if err != nil {
      t.Fatal(err)
   }
   for _, inf := range seg.inf {
      fmt.Printf("%+v\n", inf)
   }
   fmt.Printf("%+v\n", seg.key)
}


