package crypto

import (
   "bytes"
   "encoding/binary"
   "fmt"
   "github.com/89z/parse"
   "github.com/refraction-networking/utls"
   "io"
   "net"
   "net/http"
   "strings"
)

func Handshakes(pcap []byte) [][]byte {
   var hands [][]byte
   for {
      var hand []byte
      buf := NewBuffer(pcap)
      // Content Type
      junk, ok := buf.ReadBytes(0x16)
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
      pcap = pcap[len(junk):]
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

type Buffer struct {
   buf []byte
}

func NewBuffer(buf []byte) *Buffer {
   return &Buffer{buf}
}

// godocs.io/bytes#Buffer.Next
func (b *Buffer) Next(n int) ([]byte, bool) {
   if n < 0 || n > len(b.buf) {
      return nil, false
   }
   buf := b.buf[:n]
   b.buf = b.buf[n:]
   return buf, true
}

// godocs.io/bytes#Buffer.ReadBytes
func (b *Buffer) ReadBytes(delim byte) ([]byte, bool) {
   i := bytes.IndexByte(b.buf, delim)
   if i == -1 {
      return nil, false
   }
   buf := b.buf[:i+1]
   b.buf = b.buf[i+1:]
   return buf, true
}

// godocs.io/golang.org/x/crypto/cryptobyte#String.ReadUint16LengthPrefixed
func (b *Buffer) ReadUint16LengthPrefixed() ([]byte, []byte, bool) {
   low := 2
   if len(b.buf) < low {
      return nil, nil, false
   }
   high := low + int(binary.BigEndian.Uint16(b.buf))
   if len(b.buf) < high {
      return nil, nil, false
   }
   pre, buf := b.buf[:low], b.buf[low:high]
   b.buf = b.buf[high:]
   return pre, buf, true
}

// github.com/golang/go/issues/49227
func (b *Buffer) ReadUint32LengthPrefixed() ([]byte, []byte, bool) {
   low := 4
   if len(b.buf) < low {
      return nil, nil, false
   }
   high := low + int(binary.BigEndian.Uint32(b.buf))
   if len(b.buf) < high {
      return nil, nil, false
   }
   pre, buf := b.buf[:low], b.buf[low:high]
   b.buf = b.buf[high:]
   return pre, buf, true
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
      return nil, parse.Invalid{str}
   }
   var version uint16
   _, err := fmt.Sscan(tokens[0], &version)
   if err != nil {
      return nil, err
   }
   hello := ClientHello{
      new(tls.ClientHelloSpec), version,
   }
   // build CipherSuites
   cipherKeys := strings.Split(tokens[1], "-")
   for _, cipherKey := range cipherKeys {
      var cipher uint16
      _, err := fmt.Sscan(cipherKey, &cipher)
      if err != nil {
         return nil, err
      }
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
            _, err := fmt.Sscan(curveKey, &curve)
            if err != nil {
               return nil, err
            }
            curves = append(curves, curve)
         }
         ext = &tls.SupportedCurvesExtension{curves}
      case "11":
         var points []uint8
         pointKeys := strings.Split(tokens[4], "-")
         for _, pointKey := range pointKeys {
            var point uint8
            _, err := fmt.Sscan(pointKey, &point)
            if err != nil {
               return nil, err
            }
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
   for key, val := range h.Extensions {
      if key >= 1 {
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
