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
   f, err := os.Open("PCAPdroid_22_Oct_15_19_28.pcap")
   if err != nil {
      panic(err)
   }
   defer f.Close()
   hands, err := pcap.Handshakes(f)
   if err != nil {
      panic(err)
   }
   fp := tls.Fingerprinter{AllowBluntMimicry: true}
   spec, err := fp.FingerprintClientHello(hands[0])
   if err != nil {
      panic(err)
   }
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
