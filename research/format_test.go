package format

import (
   "io"
   "net/http"
   "testing"
)

func TestProgress(t *testing.T) {
   res, err := http.Get("https://speedtest.lax.hivelocity.net/100mb.file")
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   pro := NewProgress(io.Discard, res.ContentLength)
   if _, err := io.Copy(pro, res.Body); err != nil {
      t.Fatal(err)
   }
}
