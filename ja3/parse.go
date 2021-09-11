package ja3

import (
   "crypto/sha256"
   "encoding/hex"
   "fmt"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "strconv"
   "strings"
)

// extMap maps extension values to the TLSExtension object associated with the
// number. Some values are not put in here because they must be applied in a
// special way. For example, "10" is the SupportedCurves extension which is also
// used to calculate the JA3 signature. These JA3-dependent values are applied
// after the instantiation of the map.
var extMap = map[string]tls.TLSExtension{
   "0": &tls.SNIExtension{},
   "5": &tls.StatusRequestExtension{},
   "13": &tls.SignatureAlgorithmsExtension{
      SupportedSignatureAlgorithms: []tls.SignatureScheme{
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
      Modes: []uint8{tls.PskModeDHE},
   },
   "49": &tls.GenericExtension{Id: 49}, // post_handshake_auth
   "51": &tls.KeyShareExtension{
      KeyShares: []tls.KeyShare{},
   },
   "13172": &tls.NPNExtension{},
   "65281": &tls.RenegotiationInfoExtension{
      Renegotiation: tls.RenegotiateOnceAsClient,
   },
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

func Fingerprint(hello string) (*tls.ClientHelloSpec, error) {
   data, err := hex.DecodeString(hello)
   if err != nil {
      return nil, err
   }
   f := tls.Fingerprinter{AllowBluntMimicry: true}
   return f.FingerprintClientHello(data)
}

// Parse creates a ClientHelloSpec based on a JA3 string.
func Parse(ja3 string) (*tls.ClientHelloSpec, error) {
   tokens := strings.Split(ja3, ",")
   // build CipherSuites
   ciphers := strings.Split(tokens[1], "-")
   var suites []uint16
   for _, c := range ciphers {
      cid, err := strconv.ParseUint(c, 10, 16)
      if err != nil {
         return nil, err
      }
      suites = append(suites, uint16(cid))
   }
   // set extension 10
   curves := strings.Split(tokens[3], "-")
   if len(curves) == 1 && curves[0] == "" {
      curves = []string{}
   }
   var targetCurves []tls.CurveID
   for _, c := range curves {
      cid, err := strconv.ParseUint(c, 10, 16)
      if err != nil {
         return nil, err
      }
      targetCurves = append(targetCurves, tls.CurveID(cid))
   }
   extMap["10"] = &tls.SupportedCurvesExtension{Curves: targetCurves}
   // set extension 11
   pointFormats := strings.Split(tokens[4], "-")
   if len(pointFormats) == 1 && pointFormats[0] == "" {
      pointFormats = []string{}
   }
   var targetPointFormats []byte
   for _, p := range pointFormats {
      pid, err := strconv.ParseUint(p, 10, 8)
      if err != nil {
         return nil, err
      }
      targetPointFormats = append(targetPointFormats, byte(pid))
   }
   extMap["11"] = &tls.SupportedPointsExtension{
      SupportedPoints: targetPointFormats,
   }
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
   var exts []tls.TLSExtension
   for _, ext := range extensions {
      te, ok := extMap[ext]
      if !ok {
         return nil, fmt.Errorf("extension does not exist %q", ext)
      }
      exts = append(exts, te)
   }
   // return
   return &tls.ClientHelloSpec{
      CipherSuites: suites,
      CompressionMethods: []byte{0},
      Extensions: exts,
      GetSessionID: sha256.Sum256,
   }, nil
}
