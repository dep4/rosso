package crypto

import (
   "fmt"
   "github.com/89z/format"
   "github.com/refraction-networking/utls"
   "strings"
)

func ParseJA3(str string) (*tls.ClientHelloSpec, error) {
   tokens := strings.Split(str, ",")
   if tLen := len(tokens); tLen <= 4 {
      return nil, format.InvalidSlice{4, tLen}
   }
   // uTLS does not support 0x0 as min version
   hello := tls.ClientHelloSpec{TLSVersMin: tls.VersionTLS10}
   // TLSVersMin is the record version, TLSVersMax is the handshake version
   _, err := fmt.Sscan(tokens[0], &hello.TLSVersMax)
   if err != nil {
      return nil, err
   }
   // build CipherSuites
   cipherKeys := strings.Split(tokens[1], "-")
   for _, cipherKey := range cipherKeys {
      var cipher uint16
      _, err := fmt.Sscan(cipherKey, &cipher)
      if err != nil {
         return nil, err
      }
      hello.CipherSuites = append(hello.CipherSuites, cipher)
   }
   // build extenions list
   extKeys := strings.Split(tokens[2], "-")
   for _, extKey := range extKeys {
      var ext tls.TLSExtension
      switch extKey {
      case "0":
         ext = &tls.SNIExtension{}
      case "10":
         var curves []tls.CurveID
         curveKeys := strings.Split(tokens[3], "-")
         for _, curveKey := range curveKeys {
            var curve tls.CurveID
            _, err := fmt.Sscan(curveKey, &curve)
            if err != nil {
               return nil, err
            }
            curves = append(curves, curve)
         }
         ext = &tls.SupportedCurvesExtension{curves}
      case "11":
         var points []uint8
         pointKeys := strings.Split(tokens[4], "-")
         for _, pointKey := range pointKeys {
            var point uint8
            _, err := fmt.Sscan(pointKey, &point)
            if err != nil {
               return nil, err
            }
            points = append(points, point)
         }
         ext = &tls.SupportedPointsExtension{points}
      case "13":
         // this cant be empty, so just use the Go default
         ext = &tls.SignatureAlgorithmsExtension{
            []tls.SignatureScheme{
               0x804, 0x403, 0x807, 0x805, 0x806, 0x401,
               0x501, 0x601, 0x503, 0x603, 0x201, 0x203,
            },
         }
      case "16":
         // if we leave this empty, it will fail on any HTTP/2 servers
         ext = &tls.ALPNExtension{
            []string{"http/1.1"},
         }
      case "23":
         ext = &tls.UtlsExtendedMasterSecretExtension{}
      case "35":
         ext = &tls.SessionTicketExtension{}
      case "65281":
         ext = &tls.RenegotiationInfoExtension{}
      }
      hello.Extensions = append(hello.Extensions, ext)
   }
   return &hello, nil
}
