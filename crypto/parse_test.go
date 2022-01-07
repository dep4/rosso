package crypto

import (
   "net/http"
   "net/url"
   "strings"
   "testing"
)

func TestParse(t *testing.T) {
   hello, err := ParseJA3(AndroidJA3)
   if err != nil {
      t.Fatal(err)
   }
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "EncryptedPasswd": {passwd},
      "sdk_version": {"20"}, // Newer versions fail.
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth", strings.NewReader(val),
   )
   if err != nil {
      t.Fatal(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := Transport(hello).RoundTrip(req)
   if err != nil {
      t.Fatal(err)
   }
   defer res.Body.Close()
   if res.StatusCode != http.StatusOK {
      t.Fatal(res.Status)
   }
}
