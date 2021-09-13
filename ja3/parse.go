package ja3

import (
   "crypto/sha256"
   "fmt"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "strconv"
   "strings"
)

// Default values can be deleted, but otherwise dont delete or change anything,
// to avoid making a mistake. Using a function prevents modification.
func HelloGolang() *tls.ClientHelloSpec {
   return &tls.ClientHelloSpec{
      CipherSuites:[]uint16{
         0xc02b, 0xc02f, 0xc02c, 0xc030, 0xcca9, 0xcca8, 0xc009, 0xc013, 0xc00a,
         0xc014, 0x9c, 0x9d, 0x2f, 0x35, 0xc012, 0xa, 0x1301, 0x1302, 0x1303,
      },
      CompressionMethods:[]uint8{0x0},
      Extensions:[]tls.TLSExtension{
         &tls.SNIExtension{},
         &tls.StatusRequestExtension{},
         &tls.SupportedCurvesExtension{
            Curves:[]tls.CurveID{0x1d, 0x17, 0x18, 0x19},
         },
         &tls.SupportedPointsExtension{
            SupportedPoints:[]uint8{0x0},
         },
         &tls.SignatureAlgorithmsExtension{
            SupportedSignatureAlgorithms:[]tls.SignatureScheme{
               0x804, 0x403, 0x807, 0x805, 0x806, 0x401, 0x501, 0x601, 0x503,
               0x603, 0x201, 0x203,
            },
         },
         &tls.RenegotiationInfoExtension{Renegotiation:1},
         &tls.ALPNExtension{
            AlpnProtocols:[]string{"h2", "http/1.1"},
         },
         &tls.SCTExtension{},
         &tls.SupportedVersionsExtension{
            Versions:[]uint16{0x304, 0x303, 0x302, 0x301},
         },
         &tls.KeyShareExtension{
            KeyShares:[]tls.KeyShare{
               tls.KeyShare{Group:0x1d},
            },
         },
      },
   }
}

// NewTransport creates an http.Transport which mocks the given JA3 signature
// when HTTPS is used.
func NewTransport(spec *tls.ClientHelloSpec) *http.Transport {
   return &http.Transport{
      DialTLS: func(network, addr string) (net.Conn, error) {
         dialConn, err := net.Dial(network, addr)
         if err != nil {
            return nil, err
         }
         config := &tls.Config{
            ServerName: strings.Split(addr, ":")[0],
         }
         uTLSConn := tls.UClient(dialConn, config, tls.HelloCustom)
         if err := uTLSConn.ApplyPreset(spec); err != nil {
            return nil, err
         }
         if err := uTLSConn.Handshake(); err != nil {
            return nil, err
         }
         return uTLSConn, nil
      },
   }
}

// Parse creates a ClientHelloSpec based on a JA3 string.
// iana.org/assignments/tls-extensiontype-values/tls-extensiontype-values.xhtml
func Parse(ja3 string) (*tls.ClientHelloSpec, error) {
   // This must be local, to prevent mutation.
   exts := make(map[string]tls.TLSExtension)
   tokens := strings.Split(ja3, ",")
   spec := &tls.ClientHelloSpec{
      CompressionMethods: []byte{0}, GetSessionID: sha256.Sum256,
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
   // set extension 43
   vid64, err := strconv.ParseUint(tokens[0], 10, 16)
   if err != nil {
      return nil, err
   }
   exts["43"] = &tls.SupportedVersionsExtension{
      version(uint16(vid64)),
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
   return spec, nil
}

func version(min uint16) []uint16 {
   vs := []uint16{772, 771, 770, 769, 768}
   for k, v := range vs {
      if v == min {
         return vs[:k+1]
      }
   }
   return nil
}
