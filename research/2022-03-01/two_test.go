package hls

import (
   "fmt"
   "net/http"
   "testing"
)

func TestTwo(t *testing.T) {
   href, err := getHref()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(href)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   mass, err := two(res)
   if err != nil {
      t.Fatal(err)
   }
   for _, mas := range mass {
      fmt.Printf("%+v\n", mas)
   }
}
