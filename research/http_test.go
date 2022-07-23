package http

import (
   "fmt"
   "net/http"
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
   req2.Write(os.Stdout)
   res, err := new(http.Client).Do(req2)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   fmt.Println(res.Status)
}
