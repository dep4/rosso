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
   // FingerprintClientHello does bounds checking, so we dont need to worry
   // about it in this function.
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
