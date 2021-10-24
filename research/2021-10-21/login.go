package main

import (
   "bufio"
   "fmt"
   "github.com/89z/parse/pcap"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
   "time"
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
   for _, hand := range hands {
      fmt.Println(hand)
      spec, err := fp.FingerprintClientHello(hand)
      if err != nil {
         fmt.Println(err)
         continue
      }
      var exts []tls.TLSExtension
      for _, ext := range spec.Extensions {
         switch v := ext.(type) {
         case *tls.ALPNExtension:
            exts = append(exts, &tls.ALPNExtension{
               []string{"http/1.1"},
            })
         case *tls.GenericExtension:
            if v.Id != 0x7550 {
               exts = append(exts, ext)
            }
         default:
            exts = append(exts, ext)
         }
      }
      spec.Extensions = exts
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
      if res.StatusCode == http.StatusOK {
         break
      }
      time.Sleep(time.Second)
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
