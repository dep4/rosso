package m3u

import (
   "encoding/json"
   "fmt"
   "net/http"
   "path"
   "strings"
   "testing"
)

func TestM3U(t *testing.T) {
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
