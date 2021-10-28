package tls

import (
   "crypto/sha256"
   "fmt"
   "github.com/89z/parse/tls"
   "strconv"
   "strings"
   utls "github.com/refraction-networking/utls"
)

func Parse(ja3 string) (*tls.ClientHello, error) {
   // This must be local, to prevent mutation.
   exts := make(map[string]utls.TLSExtension)
   tokens := strings.Split(ja3, ",")
   spec := &utls.ClientHelloSpec{
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
   exts["0"] = &utls.SNIExtension{}
   // set extension 5
   exts["5"] = &utls.StatusRequestExtension{}
   // set extension 10
   curves := tokens[3]
   var ids []utls.CurveID
   for _, c := range strings.Split(curves, "-") {
      if c != "" {
         cid, err := strconv.ParseUint(c, 10, 16)
         if err != nil {
            return nil, err
         }
         ids = append(ids, utls.CurveID(cid))
      }
   }
   exts["10"] = &utls.SupportedCurvesExtension{ids}
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
   exts["11"] = &utls.SupportedPointsExtension{pids}
   // set extension 13
   exts["13"] = &utls.SignatureAlgorithmsExtension{
      []utls.SignatureScheme{
         utls.ECDSAWithP256AndSHA256,
         utls.PSSWithSHA256,
         utls.PKCS1WithSHA256,
         utls.ECDSAWithP384AndSHA384,
         utls.PSSWithSHA384,
         utls.PKCS1WithSHA384,
         utls.PSSWithSHA512,
         utls.PKCS1WithSHA512,
         utls.PKCS1WithSHA1,
      },
   }
   // set extension 16
   exts["16"] = &utls.ALPNExtension{
      []string{"http/1.1"},
   }
   // set extension 18
   exts["18"] = &utls.SCTExtension{}
   // set extension 21
   exts["21"] = &utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle}
   // set extension 22
   exts["22"] = &utls.GenericExtension{Id: 22} // encrypt_then_mac
   // set extension 23
   exts["23"] = &utls.UtlsExtendedMasterSecretExtension{}
   // set extension 27
   exts["27"] = &utls.FakeCertCompressionAlgsExtension{}
   // set extension 28
   exts["28"] = &utls.FakeRecordSizeLimitExtension{}
   // set extension 35
   exts["35"] = &utls.SessionTicketExtension{}
   // Set extension 43. JA3 does not specify what these should be, so just use
   // Golang default.
   exts["43"] = &utls.SupportedVersionsExtension{
      []uint16{0x304, 0x303, 0x302, 0x301},
   }
   // set extension 45
   exts["45"] = &utls.PSKKeyExchangeModesExtension{
      []uint8{utls.PskModeDHE},
   }
   // set extension 49
   exts["49"] = &utls.GenericExtension{Id: 49} // post_handshake_auth
   // set extension 50
   exts["50"] = &utls.GenericExtension{Id: 50} // signature_algorithms_cert
   // set extension 51
   exts["51"] = &utls.KeyShareExtension{
      KeyShares:[]utls.KeyShare{
         utls.KeyShare{Group:0x1d},
      },
   }
   // set extension 13172
   exts["13172"] = &utls.NPNExtension{}
   // set extension 65281
   exts["65281"] = &utls.RenegotiationInfoExtension{Renegotiation:1}
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
   return &tls.ClientHello{
      spec, uint16(ver),
   }, nil
}
