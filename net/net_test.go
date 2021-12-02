package net

import (
   "net/http/httputil"
   "os"
   "strings"
   "testing"
)

const details = `GET /fdfe/details?doc=com.instagram.android HTTP/1.1
Host: android.clients.google.com
Authorization: Bearer ya29.a0ARrdaM9rMbLUSP6wDIQGuLUH7Ej7vodCwekHOJx8J_JRD2k1...
X-DFE-Device-ID: 3a1c36f387b...

`

func TestRequest(t *testing.T) {
   req, err := ReadRequest(strings.NewReader(details))
   if err != nil {
      t.Fatal(err)
   }
   buf, err := httputil.DumpRequest(req, false)
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(buf)
}
