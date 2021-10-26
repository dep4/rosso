package main

import (
   "github.com/89z/parse/tls"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

const pass = "16030100bb010000b703034420d198e7852decbc117dc7f90550b98f2d643c954bf3361ddaf127ff921b04000024c02bc02ccca9c02fc030cca8009e009fc009c00ac013c01400330039009c009d002f00350100006aff0100010000000022002000001d636c69656e7473657276696365732e676f6f676c65617069732e636f6d0017000000230000000d0016001406010603050105030401040303010303020102030010000b000908687474702f312e31000b00020100000a000400020017"

func main() {
   hand, err := tls.Decode(pass)
   if err != nil {
      panic(err)
   }
   spec, err := hand.ClientHello()
   if err != nil {
      panic(err)
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
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   res, err := tls.NewTransport(spec).RoundTrip(req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   dum, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(dum)
}
