package hls

import (
   "fmt"
   "net/http"
   "testing"
)

func TestTwo(t *testing.T) {
   mas, err := twoHref()
   if err != nil {
      t.Fatal(err)
   }
   res, err := http.Get(mas.stream[0].URI)
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

