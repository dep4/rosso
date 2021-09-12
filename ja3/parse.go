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

var HelloGolang = &tls.ClientHelloSpec{
   CipherSuites: []uint16{
      0xC02B, // TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256
      0xC02F, // TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
      0xC02C, // TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
      0xC030, // TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
      0xCCA9, // TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256
      0xCCA8, // TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256
      0xC009, // TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA
      0xC013, // TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA
      0xC00A, // TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA
      0xC014, // TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA
      0x009C, // TLS_RSA_WITH_AES_128_GCM_SHA256
      0x009D, // TLS_RSA_WITH_AES_256_GCM_SHA384
      0x002F, // TLS_RSA_WITH_AES_128_CBC_SHA
      0x0035, // TLS_RSA_WITH_AES_256_CBC_SHA
      0xC012, // TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA
      0x000A, // TLS_RSA_WITH_3DES_EDE_CBC_SHA
      0x1301, // TLS_AES_128_GCM_SHA256
      0x1302, // TLS_AES_256_GCM_SHA384
      0x1303, // TLS_CHACHA20_POLY1305_SHA256
   },
   CompressionMethods: []uint8{0x00},
   Extensions: []tls.TLSExtension{
      &tls.SNIExtension{},
      &tls.StatusRequestExtension{},
      &tls.SupportedCurvesExtension{
         []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
            tls.CurveP384,
            tls.CurveP521,
         },
      },
      &tls.SupportedPointsExtension{
         SupportedPoints: []uint8{0x00},
      },
      &tls.SignatureAlgorithmsExtension{
         SupportedSignatureAlgorithms: []tls.SignatureScheme{
            0x804,
            0x403,
            0x807,
            0x805,
            0x806,
            0x401,
            0x501,
            0x601,
            0x503,
            0x603,
            0x201,
            0x203,
         },
      },
      &tls.RenegotiationInfoExtension{Renegotiation: 1},
      &tls.ALPNExtension{
         AlpnProtocols: []string{"h2", "http/1.1"},
      },
      &tls.SCTExtension{},
      &tls.SupportedVersionsExtension{
         Versions: []uint16{0x304, 0x303, 0x302, 0x301},
      },
      &tls.KeyShareExtension{
         KeyShares: []tls.KeyShare{
            tls.KeyShare{Group: 0x1D},
         },
      },
   },
}

// extMap maps extension values to the TLSExtension object associated with the
// number. Some values are not put in here because they must be applied in a
// special way. For example, "10" is the SupportedCurves extension which is also
// used to calculate the JA3 signature. These JA3-dependent values are applied
// after the instantiation of the map.
var extMap = map[string]tls.TLSExtension{
   "0": &tls.SNIExtension{},
   "5": &tls.StatusRequestExtension{},
   "13": &tls.SignatureAlgorithmsExtension{
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
   },
   "16": &tls.ALPNExtension{
      []string{"http/1.1"},
   },
   "18": &tls.SCTExtension{},
   "21": &tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
   "22": &tls.GenericExtension{Id: 22}, // encrypt_then_mac
   "23": &tls.UtlsExtendedMasterSecretExtension{},
   "27": &tls.FakeCertCompressionAlgsExtension{},
   "28": &tls.FakeRecordSizeLimitExtension{},
   "35": &tls.SessionTicketExtension{},
   "44": &tls.CookieExtension{},
   "45": &tls.PSKKeyExchangeModesExtension{
      []uint8{tls.PskModeDHE},
   },
   "49": &tls.GenericExtension{Id: 49}, // post_handshake_auth
   "51": &tls.KeyShareExtension{},
   "13172": &tls.NPNExtension{},
   "65281": &tls.RenegotiationInfoExtension{tls.RenegotiateOnceAsClient},
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
func Parse(ja3 string) (*tls.ClientHelloSpec, error) {
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
   extMap["10"] = &tls.SupportedCurvesExtension{ids}
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
   extMap["11"] = &tls.SupportedPointsExtension{pids}
   // set extension 43
   vid64, err := strconv.ParseUint(tokens[0], 10, 16)
   if err != nil {
      return nil, err
   }
   extMap["43"] = &tls.SupportedVersionsExtension{
      []uint16{
         uint16(vid64),
      },
   }
   // build extenions list
   extensions := strings.Split(tokens[2], "-")
   for _, ext := range extensions {
      te, ok := extMap[ext]
      if !ok {
         return nil, fmt.Errorf("extension does not exist %q", ext)
      }
      spec.Extensions = append(spec.Extensions, te)
   }
   // return
   return spec, nil
}
