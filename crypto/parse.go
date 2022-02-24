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
         // Android API 24
         ext = &tls.SNIExtension{}
      case "5":
         // Android API 26
         ext = &tls.StatusRequestExtension{}
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
         ext = &tls.SignatureAlgorithmsExtension{
            SupportedSignatureAlgorithms: []tls.SignatureScheme{
               // Android API 24
               tls.ECDSAWithP256AndSHA256,
               // httpbin.org
               tls.PKCS1WithSHA256,
            },
         }
      case "16":
         ext = &tls.ALPNExtension{
            AlpnProtocols: []string{
               // Android API 24
               "http/1.1",
            },
         }
      case "23":
         // Android API 24
         ext = &tls.UtlsExtendedMasterSecretExtension{}
      case "43":
         // Android API 29
         ext = &tls.SupportedVersionsExtension{
            Versions: []uint16{tls.VersionTLS12},
         }
      case "45":
         // Android API 29
         ext = &tls.PSKKeyExchangeModesExtension{
            Modes: []uint8{tls.PskModeDHE},
         }
      case "65281":
         // Android API 24
         ext = &tls.RenegotiationInfoExtension{}
      default:
         var id uint16
         _, err := fmt.Sscan(extKey, &id)
         if err != nil {
            return nil, err
         }
         ext = &tls.GenericExtension{Id: id}
      }
      hello.Extensions = append(hello.Extensions, ext)
   }
   return &hello, nil
}
