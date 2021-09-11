package main

import (
   "encoding/hex"
   "encoding/json"
   "fmt"
   "github.com/refraction-networking/utls"
   "os"
)

func main() {
   file, err := os.Open("curl.json")
   if err != nil {
      panic(err)
   }
   defer file.Close()
   var finger struct {
      Client_Hello string
   }
   json.NewDecoder(file).Decode(&finger)
   b, err := hex.DecodeString(finger.Client_Hello)
   if err != nil {
      panic(err)
   }
   fp := tls.Fingerprinter{AllowBluntMimicry: true}
   spec, err := fp.FingerprintClientHello(b)
   if err != nil {
      panic(err)
   }
   for _, ext := range spec.Extensions {
      fmt.Printf("%#v\n", ext)
   }
}
