package tls

import (
   "crypto/sha256"
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
   // This must be local, to prevent mutation.
   exts := make(map[string]tls.TLSExtension)
   tokens := strings.Split(ja3, ",")
   spec := &tls.ClientHelloSpec{
      CompressionMethods: []byte{0},
      GetSessionID: sha256.Sum256,
   }
   // build CipherSuites
   ciphers := tokens[1]
   for _, c := range strings.Split(ciphers, "-") {
      cid, err := strconv.ParseUint(c, 10, 16)
      if err != nil {
         return nil, err
      }
      spec.CipherSuites = append(spec.CipherSuites, uint16(cid))
   }
   // set extension 0
   exts["0"] = &tls.SNIExtension{}
   // set extension 5
   exts["5"] = &tls.StatusRequestExtension{}
   // set extension 10
   curves := tokens[3]
   var ids []tls.CurveID
   for _, c := range strings.Split(curves, "-") {
      if c != "" {
         cid, err := strconv.ParseUint(c, 10, 16)
         if err != nil {
            return nil, err
         }
         ids = append(ids, tls.CurveID(cid))
      }
   }
   exts["10"] = &tls.SupportedCurvesExtension{ids}
   // set extension 11
   pointFmts := tokens[4]
   var pids []byte
   for _, p := range strings.Split(pointFmts, "-") {
      if p != "" {
         pid, err := strconv.ParseUint(p, 10, 8)
         if err != nil {
            return nil, err
         }
         pids = append(pids, byte(pid))
      }
   }
   exts["11"] = &tls.SupportedPointsExtension{pids}
   // set extension 13
   exts["13"] = &tls.SignatureAlgorithmsExtension{
      []tls.SignatureScheme{
         tls.ECDSAWithP256AndSHA256,
         tls.PSSWithSHA256,
         tls.PKCS1WithSHA256,
         tls.ECDSAWithP384AndSHA384,
         tls.PSSWithSHA384,
         tls.PKCS1WithSHA384,
         tls.PSSWithSHA512,
         tls.PKCS1WithSHA512,
         tls.PKCS1WithSHA1,
      },
   }
   // set extension 16
   exts["16"] = &tls.ALPNExtension{
      []string{"http/1.1"},
   }
   // set extension 18
   exts["18"] = &tls.SCTExtension{}
   // set extension 21
   exts["21"] = &tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle}
   // set extension 22
   exts["22"] = &tls.GenericExtension{Id: 22} // encrypt_then_mac
   // set extension 23
   exts["23"] = &tls.UtlsExtendedMasterSecretExtension{}
   // set extension 27
   exts["27"] = &tls.FakeCertCompressionAlgsExtension{}
   // set extension 28
   exts["28"] = &tls.FakeRecordSizeLimitExtension{}
   // set extension 35
   exts["35"] = &tls.SessionTicketExtension{}
   // Set extension 43. JA3 does not specify what these should be, so just use
   // Golang default.
   exts["43"] = &tls.SupportedVersionsExtension{
      []uint16{0x304, 0x303, 0x302, 0x301},
   }
   // set extension 45
   exts["45"] = &tls.PSKKeyExchangeModesExtension{
      []uint8{tls.PskModeDHE},
   }
   // set extension 49
   exts["49"] = &tls.GenericExtension{Id: 49} // post_handshake_auth
   // set extension 50
   exts["50"] = &tls.GenericExtension{Id: 50} // signature_algorithms_cert
   // set extension 51
   exts["51"] = &tls.KeyShareExtension{
      KeyShares:[]tls.KeyShare{
         tls.KeyShare{Group:0x1d},
      },
   }
   // set extension 13172
   exts["13172"] = &tls.NPNExtension{}
   // set extension 65281
   exts["65281"] = &tls.RenegotiationInfoExtension{Renegotiation:1}
   // build extenions list
   extensions := strings.Split(tokens[2], "-")
   for _, ext := range extensions {
      if ext == "10" && curves == "" {
         return nil, fmt.Errorf("SSLExtension %q EllipticCurve %q", ext, curves)
      }
      te, ok := exts[ext]
      if !ok {
         return nil, fmt.Errorf("extension does not exist %q", ext)
      }
      spec.Extensions = append(spec.Extensions, te)
   }
   // return
   ver, err := strconv.Atoi(tokens[0])
   if err != nil {
      return nil, err
   }
   return &ClientHello{
      spec, uint16(ver),
   }, nil
}

func extensionType(ext tls.TLSExtension) (uint16, error) {
   data, err := io.ReadAll(ext)
   if err != nil {
      return 0, err
   }
   return binary.BigEndian.Uint16(data), nil
}
