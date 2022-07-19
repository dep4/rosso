// TLS and JA3 parsers
package crypto

import (
   "fmt"
   "github.com/refraction-networking/utls"
   "strings"
)

func Format_JA3(spec *tls.ClientHelloSpec) (string, error) {
   buf := new(strings.Builder)
   // TLSVersMin is the record version, TLSVersMax is the handshake version
   fmt.Fprint(buf, spec.TLSVersMax)
   // Cipher Suites
   buf.WriteByte(',')
   for key, val := range spec.CipherSuites {
      if key >= 1 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, val)
   }
   // Extensions
   buf.WriteByte(',')
   var (
      curves []tls.CurveID
      points []uint8
   )
   for key, val := range spec.Extensions {
      switch ext := val.(type) {
      case *tls.SupportedCurvesExtension:
         curves = ext.Curves
      case *tls.SupportedPointsExtension:
         points = ext.SupportedPoints
      }
      typ, err := extension_type(val)
      if err != nil {
         return "", err
      }
      if key >= 1 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, typ)
   }
   // Elliptic curves
   buf.WriteByte(',')
   for key, val := range curves {
      if key >= 1 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, val)
   }
   // ECPF
   buf.WriteByte(',')
   for key, val := range points {
      if key >= 1 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, val)
   }
   return buf.String(), nil
}
