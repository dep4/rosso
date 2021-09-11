package main

import (
   "bufio"
   "encoding/hex"
   "encoding/json"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
)

func main() {
   res, err := http.Get("https://client.tlsfingerprint.io:8443")
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   var fp struct {
      Client_Hello string
   }
   json.NewDecoder(res.Body).Decode(&fp)
   data, err := hex.DecodeString(fp.Client_Hello)
   if err != nil {
      panic(err)
   }
   spec, err := new(tls.Fingerprinter).FingerprintClientHello(data)
   if err != nil {
      panic(err)
   }
   req, err := http.NewRequest("HEAD", "https://example.com", nil)
   if err != nil {
      panic(err)
   }
   conn, err := net.Dial("tcp", req.URL.Host + ":" + req.URL.Scheme)
   if err != nil {
      panic(err)
   }
   cfg := &tls.Config{ServerName: req.URL.Host}
   uConn := tls.UClient(conn, cfg, tls.HelloCustom)
   if err := uConn.ApplyPreset(spec); err != nil {
      panic(err)
   }
   if err := req.Write(uConn); err != nil {
      panic(err)
   }
   if _, err := http.ReadResponse(bufio.NewReader(uConn), req); err != nil {
      panic(err)
   }
}
