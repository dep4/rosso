package hls

import (
   "fmt"
   "net/http"
   "testing"
)

func twoHref() ([]master, error) {
   href, err := oneHref()
   if err != nil {
      return nil, err
   }
   res, err := http.Get(href)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   return two(res)
}

func TestTwo(t *testing.T) {
   mass, err := twoHref()
   if err != nil {
      t.Fatal(err)
   }
   for _, mas := range mass {
      fmt.Printf("%+v\n", mas)
   }
}
