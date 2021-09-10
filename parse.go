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
      AlpnProtocols: []string{"h2", "http/1.1"},
   },
   "18": &tls.SCTExtension{},
   "22": &tls.GenericExtension{Id: 22}, // encrypt_then_mac
   "21": &tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
   "23": &tls.UtlsExtendedMasterSecretExtension{},
   "27": &tls.FakeCertCompressionAlgsExtension{},
   "28": &tls.FakeRecordSizeLimitExtension{},
   "35": &tls.SessionTicketExtension{},
   "43": &tls.SupportedVersionsExtension{
      Versions: []uint16{
         tls.GREASE_PLACEHOLDER,
         tls.VersionTLS13,
         tls.VersionTLS12,
         tls.VersionTLS11,
         tls.VersionTLS10,
      },
   },
   "44": &tls.CookieExtension{},
   "45": &tls.PSKKeyExchangeModesExtension{
      Modes: []uint8{tls.PskModeDHE},
   },
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

// Parse creates a ClientHelloSpec based on a JA3 string.
func Parse(ja3 string) (*tls.ClientHelloSpec, error) {
   tokens := strings.Split(ja3, ",")
   version := tokens[0]
   ciphers := strings.Split(tokens[1], "-")
   extensions := strings.Split(tokens[2], "-")
   curves := strings.Split(tokens[3], "-")
   if len(curves) == 1 && curves[0] == "" {
      curves = []string{}
   }
   pointFormats := strings.Split(tokens[4], "-")
   if len(pointFormats) == 1 && pointFormats[0] == "" {
      pointFormats = []string{}
   }
   // parse curves
   var targetCurves []tls.CurveID
   for _, c := range curves {
      cid, err := strconv.ParseUint(c, 10, 16)
      if err != nil {
         return nil, fmt.Errorf("curve %v", err)
      }
      targetCurves = append(targetCurves, tls.CurveID(cid))
   }
   extMap["10"] = &tls.SupportedCurvesExtension{Curves: targetCurves}
   // parse point formats
   var targetPointFormats []byte
   for _, p := range pointFormats {
      pid, err := strconv.ParseUint(p, 10, 8)
      if err != nil {
         return nil, fmt.Errorf("pointFormat %v", err)
      }
      targetPointFormats = append(targetPointFormats, byte(pid))
   }
   extMap["11"] = &tls.SupportedPointsExtension{
      SupportedPoints: targetPointFormats,
   }
   // build extenions list
   var exts []tls.TLSExtension
   for _, ext := range extensions {
      te, ok := extMap[ext]
      if !ok {
         return nil, fmt.Errorf("extension does not exist %q", ext)
      }
      exts = append(exts, te)
   }
   // build SSLVersion
   vid64, err := strconv.ParseUint(version, 10, 16)
   if err != nil {
      return nil, fmt.Errorf("version %v", err)
   }
   vid := uint16(vid64)
   // build CipherSuites
   var suites []uint16
   for _, c := range ciphers {
      cid, err := strconv.ParseUint(c, 10, 16)
      if err != nil {
         return nil, fmt.Errorf("cipher %v", err)
      }
      suites = append(suites, uint16(cid))
   }
   return &tls.ClientHelloSpec{
      CipherSuites: suites,
      CompressionMethods: []byte{0},
      Extensions: exts,
      GetSessionID: sha256.Sum256,
      TLSVersMax: vid,
      TLSVersMin: vid,
   }, nil
}
