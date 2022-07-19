// TLS and JA3 parsers
package crypto

import (
   "crypto/md5"
   "encoding/binary"
   "encoding/hex"
   "fmt"
   "github.com/refraction-networking/utls"
   "io"
   "net"
   "net/http"
   "strings"
)

func Parse_JA3(str string) (*tls.ClientHelloSpec, error) {
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
         for _, raw_cipher := range strings.Split(field, "-") {
            var cipher uint16
            _, err := fmt.Sscan(raw_cipher, &cipher)
            if err != nil {
               return nil, err
            }
            spec.CipherSuites = append(spec.CipherSuites, cipher)
         }
      case 2:
         extensions = field
      case 3:
         for _, raw_curve := range strings.Split(field, "-") {
            var curve tls.CurveID
            _, err := fmt.Sscan(raw_curve, &curve)
            if err != nil {
               return nil, err
            }
            info.SupportedCurves = append(info.SupportedCurves, curve)
         }
      case 4:
         for _, raw_point := range strings.Split(field, "-") {
            var point uint8
            _, err := fmt.Sscan(raw_point, &point)
            if err != nil {
               return nil, err
            }
            info.SupportedPoints = append(info.SupportedPoints, point)
         }
      }
   }
   // build extenions list
   for _, raw_ID := range strings.Split(extensions, "-") {
      var ext tls.TLSExtension
      switch raw_ID {
      case "0":
         // Android API 24
         ext = &tls.SNIExtension{}
      case "5":
         // Android API 26
         ext = &tls.StatusRequestExtension{}
      case "10":
         ext = &tls.SupportedCurvesExtension{Curves: info.SupportedCurves}
      case "11":
         ext = &tls.SupportedPointsExtension{
            SupportedPoints: info.SupportedPoints,
         }
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
         // Android API 24
         ext = &tls.ALPNExtension{
            AlpnProtocols: []string{"http/1.1"},
         }
      case "23":
         // Android API 24
         ext = &tls.UtlsExtendedMasterSecretExtension{}
      case "27":
         // Google Chrome
         ext = &tls.FakeCertCompressionAlgsExtension{
            Methods: []tls.CertCompressionAlgo{tls.CertCompressionBrotli},
         }
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
         _, err := fmt.Sscan(raw_ID, &id)
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

// cannot call pointer method RoundTrip on http.Transport
func Transport(spec *tls.ClientHelloSpec) *http.Transport {
   var tr http.Transport
   //lint:ignore SA1019 godocs.io/context
   tr.DialTLS = func(network, ref string) (net.Conn, error) {
      conn, err := net.Dial(network, ref)
      if err != nil {
         return nil, err
      }
      host, _, err := net.SplitHostPort(ref)
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
   }
   return &tr
}

// len 122, 8fcaa9e4a15f48af0a7d396e3fa5c5eb
const Android_API_24 =
   "771,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-" +
   "49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

// len 128, 9fc6ef6efc99b933c5e2d8fcf4f68955
const Android_API_25 =
   "771,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-" +
   "49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23-24-25,0"

// len 116, d8c87b9bfde38897979e41242626c2f3
const Android_API_26 =
   "771,49195-49196-52393-49199-49200-52392-49161-49162-49171-" +
   "49172-156-157-47-53,65281-0-23-35-13-5-16-11-10,29-23-24,0"

// len 143, 9b02ebd3a43b62d825e1ac605b621dc8
const Android_API_29 =
   "771,4865-4866-4867-49195-49196-52393-49199-49200-52392-49161-49162-49171-" +
   "49172-156-157-47-53,0-23-65281-10-11-35-16-5-13-51-45-43-21,29-23-24,0"

const Android_API_32 = Android_API_29

func Parse_TLS(buf []byte) (*tls.ClientHelloSpec, error) {
   // unsupported extension 0x16
   printer := tls.Fingerprinter{AllowBluntMimicry: true}
   spec, err := printer.FingerprintClientHello(buf)
   if err != nil {
      return nil, err
   }
   // If SupportedVersionsExtension is present, then TLSVersMax is set to zero.
   // In which case we need to manually read the bytes.
   if spec.TLSVersMax == 0 {
      // \x16\x03\x01\x00\xbc\x01\x00\x00\xb8\x03\x03
      spec.TLSVersMax = binary.BigEndian.Uint16(buf[9:])
   }
   return spec, nil
}

func Fingerprint(ja3 string) string {
   hash := md5.New()
   io.WriteString(hash, ja3)
   sum := hash.Sum(nil)
   return hex.EncodeToString(sum)
}

func extension_type(ext tls.TLSExtension) (uint16, error) {
   pad, ok := ext.(*tls.UtlsPaddingExtension)
   if ok {
      pad.WillPad = true
      ext = pad
   }
   buf, err := io.ReadAll(ext)
   if err != nil || len(buf) <= 1 {
      return 0, err
   }
   return binary.BigEndian.Uint16(buf), nil
}
