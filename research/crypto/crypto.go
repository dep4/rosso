package crypto

import (
   "fmt"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "strings"
)

func Transport(spec *tls.ClientHelloSpec) *http.Transport {
   return &http.Transport{
      DialTLS: func(network, addr string) (net.Conn, error) {
         conn, err := net.Dial(network, addr)
         if err != nil {
            return nil, err
         }
         host, _, err := net.SplitHostPort(addr)
         if err != nil {
            return nil, err
         }
         config := &tls.Config{ServerName: host}
         uconn := tls.UClient(conn, config, tls.HelloCustom)
         if err := uconn.ApplyPreset(spec); err != nil {
            return nil, err
         }
         if err := uconn.Handshake(); err != nil {
            return nil, err
         }
         return uconn, nil
      },
   }
}

// 9b02ebd3a43b62d825e1ac605b621dc8
const AndroidAPI29 =
   "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-" +
   "49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"

// d8c87b9bfde38897979e41242626c2f3
const AndroidAPI26 =
   "771,49195-49196-52393-49199-49200-52392-49161-49162-49171-" +
   "49172-156-157-47-53,65281-0-23-35-13-5-16-11-10,29-23-24,0"

// 9fc6ef6efc99b933c5e2d8fcf4f68955
const AndroidAPI25 =
   "771,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-" +
   "49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23-24-25,0"

// 8fcaa9e4a15f48af0a7d396e3fa5c5eb
const AndroidAPI24 =
   "771,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-" +
   "49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

func ParseJA3(str string) (*tls.ClientHelloSpec, error) {
   tokens := strings.Split(str, ",")
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
         ext = &tls.SupportedVersionsExtension{
            // Android API 29
            Versions: []uint16{tls.VersionTLS12},
         }
      case "45":
         ext = &tls.PSKKeyExchangeModesExtension{
            // Android API 29
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
