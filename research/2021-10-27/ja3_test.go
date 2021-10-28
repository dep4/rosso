package tls

import (
   "fmt"
   "github.com/89z/parse/tls"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
   "testing"
)

func TestParse(t *testing.T) {
   hello, err := parse(tls.Android)
   if err != nil {
      t.Fatal(err)
   }
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "sdk_version": {"17"},
      "EncryptedPasswd": {encryptedPasswd},
   }
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth",
      strings.NewReader(val.Encode()),
   )
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   fmt.Println("RoundTrip")
   res, err := tls.NewTransport(hello.ClientHelloSpec).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   fmt.Println("DumpResponse")
   dum, err := httputil.DumpResponse(res, true)
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(dum)
}
