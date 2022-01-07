package crypto

import (
   "fmt"
   "net/http"
   "net/url"
   "strings"
   "testing"
   "time"
)

var hellos = []string{
   AndroidJA3,
   "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0",
}

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
