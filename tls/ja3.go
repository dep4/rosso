package tls

import (
   "encoding/binary"
   "fmt"
   "github.com/refraction-networking/utls"
   "io"
   "strconv"
   "strings"
)

const Android = "769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

func Marshal(hello *ClientHello) (string, error) {
   buf := new(strings.Builder)
   // Version
   fmt.Fprint(buf, hello.Version)
   // Cipher Suites
   buf.WriteByte(',')
   for key, val := range hello.CipherSuites {
      if key > 0 {
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
   for key, val := range hello.Extensions {
      if key > 0 {
         buf.WriteByte('-')
      }
      typ, err := extensionType(val)
      if err != nil {
         return "", err
      }
      fmt.Fprint(buf, typ)
      switch ext := val.(type) {
      case *tls.SupportedCurvesExtension:
         curves = ext.Curves
      case *tls.SupportedPointsExtension:
         points = ext.SupportedPoints
      }
   }
   // Elliptic curves
   buf.WriteByte(',')
   for key, val := range curves {
      if key > 0 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, val)
   }
   // ECPF
   buf.WriteByte(',')
   for key, val := range points {
      if key > 0 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, val)
   }
   return buf.String(), nil
}

func Parse(ja3 string) (*ClientHello, error) {
   tokens := strings.Split(ja3, ",")
   ver, err := strconv.Atoi(tokens[0])
   if err != nil {
      return nil, err
   }
   hello := ClientHello{
      new(tls.ClientHelloSpec), uint16(ver),
   }
   // build CipherSuites
   ciphers := strings.Split(tokens[1], "-")
   for _, key := range ciphers {
      val, err := strconv.ParseUint(key, 10, 16)
      if err != nil {
         return nil, err
      }
      hello.CipherSuites = append(hello.CipherSuites, uint16(val))
   }
   // build extenions list
   extensions := strings.Split(tokens[2], "-")
   for _, key := range extensions {
      var val tls.TLSExtension
      switch key {
      case "0":
         val = &tls.SNIExtension{}
      case "10":
         var cids []tls.CurveID
         curves := strings.Split(tokens[3], "-")
         for _, key := range curves {
            val, err := strconv.ParseUint(key, 10, 16)
            if err != nil {
               return nil, err
            }
            cids = append(cids, tls.CurveID(val))
         }
         val = &tls.SupportedCurvesExtension{cids}
      case "11":
         var pids []uint8
         points := strings.Split(tokens[4], "-")
         for _, key := range points {
            val, err := strconv.ParseUint(key, 10, 8)
            if err != nil {
               return nil, err
            }
            pids = append(pids, uint8(val))
         }
         val = &tls.SupportedPointsExtension{pids}
      case "13":
         // this cant be empty, so just use the Go default
         val = &tls.SignatureAlgorithmsExtension{
            []tls.SignatureScheme{
               0x804, 0x403, 0x807, 0x805, 0x806, 0x401,
               0x501, 0x601, 0x503, 0x603, 0x201, 0x203,
            },
         }
      case "16":
         // if we leave this empty, it will fail on any HTTP/2 servers
         val = &tls.ALPNExtension{
            []string{"http/1.1"},
         }
      case "23":
         val = &tls.UtlsExtendedMasterSecretExtension{}
      case "35":
         val = &tls.SessionTicketExtension{}
      case "65281":
         val = &tls.RenegotiationInfoExtension{}
      }
      hello.Extensions = append(hello.Extensions, val)
   }
   return &hello, nil
}

func extensionType(ext tls.TLSExtension) (uint16, error) {
   data, err := io.ReadAll(ext)
   if err != nil {
      return 0, err
   }
   return binary.BigEndian.Uint16(data), nil
}
