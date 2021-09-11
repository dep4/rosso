package main

import (
   "bufio"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
)

func main() {
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
   if err := uConn.ApplyPreset(fail); err != nil {
      panic(err)
   }
   if err := req.Write(uConn); err != nil {
      panic(err)
   }
   if _, err := http.ReadResponse(bufio.NewReader(uConn), req); err != nil {
      panic(err)
   }
}
