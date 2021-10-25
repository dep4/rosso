package main

import (
   "bufio"
   "github.com/89z/parse/pcap"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

func main() {
   data, err := os.ReadFile("PCAPdroid_25_Oct_21_53_41.pcap")
   if err != nil {
      panic(err)
   }
   for _, hand := range pcap.Handshakes(data) {
      spec, err := hand.ClientHello()
      if err == nil {
         res, err := post(spec)
         if err != nil {
            panic(err)
         }
         defer res.Body.Close()
         dum, err := httputil.DumpResponse(res, true)
         if err != nil {
            panic(err)
         }
         os.Stdout.Write(dum)
         break
      }
   }
}

func post(spec *tls.ClientHelloSpec) (*http.Response, error) {
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
      return nil, err
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   tcpConn, err := net.Dial("tcp", req.URL.Host + ":" + req.URL.Scheme)
   if err != nil {
      return nil, err
   }
   config := &tls.Config{ServerName: req.URL.Host}
   tlsConn := tls.UClient(tcpConn, config, tls.HelloCustom)
   if err := tlsConn.ApplyPreset(spec); err != nil {
      return nil, err
   }
   if err := req.Write(tlsConn); err != nil {
      return nil, err
   }
   return http.ReadResponse(bufio.NewReader(tlsConn), req)
}
