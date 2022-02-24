package crypto

import (
   "fmt"
   "github.com/refraction-networking/utls"
   "strings"
)

func ParseJA3(str string) (*tls.ClientHelloSpec, error) {
   var (
      extensions string
      info tls.ClientHelloInfo
      spec tls.ClientHelloSpec
   )
   for i, field := range strings.SplitN(str, ",", 5) {
      switch i {
      case 0:
         // TLSVersMin is the record version, TLSVersMax is the handshake
         // version
         _, err := fmt.Sscan(field, &spec.TLSVersMax)
         if err != nil {
            return nil, err
         }
      case 1:
         // build CipherSuites
         for _, cipherKey := range strings.Split(field, "-") {
            var cipher uint16
            _, err := fmt.Sscan(cipherKey, &cipher)
            if err != nil {
               return nil, err
            }
            spec.CipherSuites = append(spec.CipherSuites, cipher)
         }
      case 2:
         extensions = field
      case 3:
         for _, curveKey := range strings.Split(field, "-") {
            var curve tls.CurveID
            _, err := fmt.Sscan(curveKey, &curve)
            if err != nil {
               return nil, err
            }
            info.SupportedCurves = append(info.SupportedCurves, curve)
         }
      case 4:
         for _, pointKey := range strings.Split(field, "-") {
            var point uint8
            _, err := fmt.Sscan(pointKey, &point)
            if err != nil {
               return nil, err
            }
            info.SupportedPoints = append(info.SupportedPoints, point)
         }
      }
   }
   // build extenions list
   for _, extKey := range strings.Split(extensions, "-") {
      var ext tls.TLSExtension
      switch extKey {
      case "0":
         // Android API 24
         ext = &tls.SNIExtension{}
      case "5":
         // Android API 26
         ext = &tls.StatusRequestExtension{}
      case "10":
         ext = &tls.SupportedCurvesExtension{info.SupportedCurves}
      case "11":
         ext = &tls.SupportedPointsExtension{info.SupportedPoints}
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
      spec.Extensions = append(spec.Extensions, ext)
   }
   // uTLS does not support 0x0 as min version
   spec.TLSVersMin = tls.VersionTLS10
   return &spec, nil
}
