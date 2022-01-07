package main

import (
   "fmt"
   "github.com/89z/format/crypto"
   "github.com/CUCyber/ja3transport"
   "net/http"
   "net/url"
   "strings"
)

var _ = crypto.AndroidJA3

const (
   //ja3 = crypto.AndroidJA3
   ja3 = "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"
)

func main() {
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "EncryptedPasswd": {passwd},
      "sdk_version": {"20"}, // Newer versions fail.
   }.Encode()
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth", strings.NewReader(val),
   )
   if err != nil {
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   tra, err := ja3transport.NewTransport(ja3)
   if err != nil {
      panic(err)
   }
   res, err := tra.RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   fmt.Printf("%+v\n", res)
}
