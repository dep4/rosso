package crypto

import (
   "fmt"
   "net/http"
   "net/url"
   "strings"
   "testing"
   "time"
)

var hellos = []string{AndroidAPI24, AndroidAPI25, AndroidAPI26, AndroidAPI29}

func TestParse(t *testing.T) {
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "EncryptedPasswd": {passwd},
      "sdk_version": {"20"}, // Newer versions fail.
   }.Encode()
   for _, hello := range hellos {
      spec, err := ParseJA3(hello)
      if err != nil {
         t.Fatal(err)
      }
      req, err := http.NewRequest(
         "POST", "https://android.clients.google.com/auth", strings.NewReader(val),
      )
      if err != nil {
         t.Fatal(err)
      }
      req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
      res, err := Transport(spec).RoundTrip(req)
      if err != nil {
         t.Fatal(err)
      }
      defer res.Body.Close()
      fmt.Println(res.Status, hello)
      time.Sleep(time.Second)
   }
}
