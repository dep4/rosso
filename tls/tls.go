package tls

import (
   "encoding/binary"
   "fmt"
   "github.com/89z/parse/bytes"
   "github.com/refraction-networking/utls"
   "io"
   "net"
   "net/http"
   "strings"
)

const Android = "769,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

func Handshakes(pcap []byte) [][]byte {
   var hands [][]byte
   for {
      var hand []byte
      buf := bytes.NewBuffer(pcap)
      // Content Type
      cut, ok := buf.ReadBytes(0x16)
      if ! ok {
         return hands
      }
      hand = append(hand, 0x16)
      // Version
      ver, ok := buf.Next(2)
      if ok {
         hand = append(hand, ver...)
      }
      // Length, Handshake Protocol
      pre, pro, ok := buf.ReadUint16LengthPrefixed()
      if ok {
         hand = append(hand, pre...)
         hand = append(hand, pro...)
         hands = append(hands, hand)
      }
      pcap = pcap[len(cut):]
   }
}

func NewTransport(spec *tls.ClientHelloSpec) *http.Transport {
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

func extensionType(ext tls.TLSExtension) (uint16, error) {
   data, err := io.ReadAll(ext)
   if err != nil {
      return 0, err
   }
   return binary.BigEndian.Uint16(data), nil
}

type ClientHello struct {
   *tls.ClientHelloSpec
   Version uint16
}

func ParseHandshake(data []byte) (*ClientHello, error) {
   version := binary.BigEndian.Uint16(data[1:])
   var fin tls.Fingerprinter
   spec, err := fin.FingerprintClientHello(data)
   if err != nil {
      return nil, err
   }
   return &ClientHello{spec, version}, nil
}

func ParseJA3(str string) (*ClientHello, error) {
   tokens := strings.Split(str, ",")
   if len(tokens) != 5 {
      return nil, fmt.Errorf("%q", tokens)
   }
   var version uint16
   fmt.Sscan(tokens[0], &version)
   hello := ClientHello{
      new(tls.ClientHelloSpec), version,
   }
   // build CipherSuites
   cipherKeys := strings.Split(tokens[1], "-")
   for _, cipherKey := range cipherKeys {
      var cipher uint16
      fmt.Sscan(cipherKey, &cipher)
      hello.CipherSuites = append(hello.CipherSuites, cipher)
   }
   // build extenions list
   extKeys := strings.Split(tokens[2], "-")
   for _, extKey := range extKeys {
      var ext tls.TLSExtension
      switch extKey {
      case "0":
         ext = &tls.SNIExtension{}
      case "10":
         var curves []tls.CurveID
         curveKeys := strings.Split(tokens[3], "-")
         for _, curveKey := range curveKeys {
            var curve tls.CurveID
            fmt.Sscan(curveKey, &curve)
            curves = append(curves, curve)
         }
         ext = &tls.SupportedCurvesExtension{curves}
      case "11":
         var points []uint8
         pointKeys := strings.Split(tokens[4], "-")
         for _, pointKey := range pointKeys {
            var point uint8
            fmt.Sscan(pointKey, &point)
            points = append(points, point)
         }
         ext = &tls.SupportedPointsExtension{points}
      case "13":
         // this cant be empty, so just use the Go default
         ext = &tls.SignatureAlgorithmsExtension{
            []tls.SignatureScheme{
               0x804, 0x403, 0x807, 0x805, 0x806, 0x401,
               0x501, 0x601, 0x503, 0x603, 0x201, 0x203,
            },
         }
      case "16":
         // if we leave this empty, it will fail on any HTTP/2 servers
         ext = &tls.ALPNExtension{
            []string{"http/1.1"},
         }
      case "23":
         ext = &tls.UtlsExtendedMasterSecretExtension{}
      case "35":
         ext = &tls.SessionTicketExtension{}
      case "65281":
         ext = &tls.RenegotiationInfoExtension{}
      }
      hello.Extensions = append(hello.Extensions, ext)
   }
   return &hello, nil
}

func (h ClientHello) FormatJA3() (string, error) {
   buf := new(strings.Builder)
   // Version
   fmt.Fprint(buf, h.Version)
   // Cipher Suites
   buf.WriteByte(',')
   for key, val := range h.CipherSuites {
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
   for key, val := range h.Extensions {
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
