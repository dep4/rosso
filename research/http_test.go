package http

import (
   "fmt"
   "net/http"
   "net/http/httputil"
   "os"
   "testing"
)

func Test_Clone(t *testing.T) {
   req, err := http.NewRequest("", "http://example.com/hello", nil)
   if err != nil {
      t.Fatal(err)
   }
   req2 := Clone(req)
   req2.URL.Host = "github.com"
   buf, err := httputil.DumpRequestOut(req2, false)
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(buf)
   res, err := new(http.Client).Do(req2)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   fmt.Printf("%+v\n", res)
}
