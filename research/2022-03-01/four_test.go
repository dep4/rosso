package hls

import (
   "fmt"
   "net/http"
   "testing"
)

func TestFour(t *testing.T) {
   mass, err := twoHref()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(mass[0].uri)
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
