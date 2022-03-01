package hls

import (
   "fmt"
   "net/http"
   "testing"
)

func TestThree(t *testing.T) {
   mass, err := twoHref()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(mass[0].uri)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   for _, seg := range three(res.Body) {
      fmt.Printf("%+v\n", seg)
   }
}
