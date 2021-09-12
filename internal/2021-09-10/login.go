package main

import (
   "bufio"
   "fmt"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

func main() {
   pass := os.Getenv("ENCRYPTEDPASS")
   if pass == "" {
      fmt.Println("missing pass")
      return
   }
   val := url.Values{
      "Email": {"srpen6@gmail.com"},
      "sdk_version": {"17"},
      "EncryptedPasswd": {pass},
   }
   req, err := http.NewRequest(
      "POST", "https://android.clients.google.com/auth",
      strings.NewReader(val.Encode()),
   )
   if err != nil {
      panic(err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   dReq, err := httputil.DumpRequest(req, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(append(dReq, '\n'))
   tcpConn, err := net.Dial("tcp", req.URL.Host + ":" + req.URL.Scheme)
   if err != nil {
      panic(err)
   }
   config := &tls.Config{ServerName: req.URL.Host}
   tlsConn := tls.UClient(tcpConn, config, tls.HelloCustom)
   if err := tlsConn.ApplyPreset(preset); err != nil {
      panic(err)
   }
   if err := req.Write(tlsConn); err != nil {
      panic(err)
   }
   res, err := http.ReadResponse(bufio.NewReader(tlsConn), req)
   if err != nil {
      panic(err)
   }
   defer res.Body.Close()
   dRes, err := httputil.DumpResponse(res, true)
   if err != nil {
      panic(err)
   }
   os.Stdout.Write(dRes)
}
