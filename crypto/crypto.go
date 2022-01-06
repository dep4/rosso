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

type Buffer struct {
   buf []byte
}

func NewBuffer(buf []byte) *Buffer {
   return &Buffer{buf}
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

// 8fcaa9e4a15f48af0a7d396e3fa5c5eb
const AndroidJA3 =
   "771,49195-49196-52393-49199-49200-52392-158-159-49161-49162-49171-49172-" +
   "51-57-156-157-47-53,65281-0-23-35-13-16-11-10,23,0"

func Fingerprint(ja3 string) string {
   hash := md5.New()
   io.WriteString(hash, ja3)
   sum := hash.Sum(nil)
   return hex.EncodeToString(sum)
}

func FormatJA3(spec *tls.ClientHelloSpec) (string, error) {
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


func ParseTLS(buf []byte) (*tls.ClientHelloSpec, error) {
   // unsupported extension 0x16
   fin := tls.Fingerprinter{AllowBluntMimicry: true}
   return fin.FingerprintClientHello(buf)
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

func extensionType(ext tls.TLSExtension) (uint16, error) {
   buf, err := io.ReadAll(ext)
   if err != nil || len(buf) <= 1 {
      return 0, err
   }
   return binary.BigEndian.Uint16(buf), nil
}
