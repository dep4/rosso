package format

import (
   "fmt"
   "io"
   "net/http"
   "testing"
)

func TestProgress(t *testing.T) {
   res, err := http.Get("http://speedtest.lax.hivelocity.net/10Mio.dat")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := NewProgress(res)
   io.ReadAll(pro)
}
