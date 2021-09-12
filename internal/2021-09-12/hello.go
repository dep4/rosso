package main

import (
   "encoding/hex"
   "encoding/json"
   "fmt"
   "github.com/refraction-networking/utls"
   "net/http"
)

func main() {
   r, err := http.Get("https://client.tlsfingerprint.io:8443")
   if err != nil {
      panic(err)
   }
   defer r.Body.Close()
   var f struct {
      Client_Hello string
   }
   json.NewDecoder(r.Body).Decode(&f)
   data, err := hex.DecodeString(f.Client_Hello)
   if err != nil {
      panic(err)
   }
   spec, err := new(tls.Fingerprinter).FingerprintClientHello(data)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%#v\n", spec)
   for _, ext := range spec.Extensions {
      fmt.Printf("%#v\n", ext)
   }
}
