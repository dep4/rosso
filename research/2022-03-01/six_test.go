package hls

import (
   "fmt"
   "net/http"
   "testing"
)

func TestSix(t *testing.T) {
   href, err := platformOne()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(href)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   segs, err := four(res)
   if err != nil {
      t.Fatal(err)
   }
   for _, seg := range segs {
      fmt.Printf("%+v\n", seg)
   }
}
