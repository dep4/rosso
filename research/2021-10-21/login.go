package main

import (
   "bufio"
   "encoding/hex"
   "fmt"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "net/http/httputil"
   "net/url"
   "os"
   "strings"
)

const hello = "16030101fc010001f803037a746304c5306ac31f73b40f6bc5c7a358755421e97778241081d3086a841aad20d448abaf74a40ed1f063172691561dc0188de102eb62dafe0e8f626ad458721f0024130113031302c02bc02fcca9cca8c02cc030c00ac009c013c014009c009d002f0035000a0100018b0000001d001b000018636c69656e742e746c7366696e6765727072696e742e696f00170000ff01000100000a000e000c001d00170018001901000101000b00020100002300000010000e000c02683208687474702f312e310005000501000000000033006b0069001d002019241a5129054181978fb8c10a90f0a1ce7f77eff71d66688339c985182a54710017004104edbd4e6fe01de365be0b1f5d79b4328d7d6ccde4e31a8c68e24c428708c2a908ce065a1a059ce897ee353285063f165a2b7af61540357cd8bd5fed6f850ffc1f002b00050403040303000d0018001604030503060308040805080604010501060102030201002d00020101001c00024001001500860000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

var _ = fmt.Print

func main() {
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
   // BEGIN
   data, err := hex.DecodeString(hello)
   if err != nil {
      panic(err)
   }
   fp := tls.Fingerprinter{AllowBluntMimicry: true}
   spec, err := fp.FingerprintClientHello(data)
   if err != nil {
      panic(err)
   }
   for k, v := range spec.Extensions {
      _, ok := v.(*tls.ALPNExtension)
      if ok {
         spec.Extensions[k] = &tls.ALPNExtension{
            []string{"http/1.1"},
         }
      }
   }
   // END
   if err := tlsConn.ApplyPreset(spec); err != nil {
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
