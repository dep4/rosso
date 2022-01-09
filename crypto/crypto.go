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
      switch ext := val.(type) {
      case *tls.SupportedCurvesExtension:
         curves = ext.Curves
      case *tls.SupportedPointsExtension:
         points = ext.SupportedPoints
      }
      typ, err := extensionType(val)
      if err != nil {
         return "", err
      }
      if key >= 1 {
         buf.WriteByte('-')
      }
      fmt.Fprint(buf, typ)
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
   pad, ok := ext.(*tls.UtlsPaddingExtension)
   if ok {
      pad.WillPad = true
      ext = pad
   }
   buf, err := io.ReadAll(ext)
   if err != nil || len(buf) <= 1 {
      return 0, err
   }
   return binary.BigEndian.Uint16(buf), nil
}

type Reader struct {
   buf []byte
}

func NewReader(buf []byte) *Reader {
   return &Reader{buf}
}

// github.com/golang/go/issues/49227
func (r *Reader) ReadUint32LengthPrefixed() ([]byte, []byte, bool) {
   low := 4
   if len(r.buf) < low {
      return nil, nil, false
   }
   high := low + int(binary.BigEndian.Uint32(r.buf))
   if len(r.buf) < high {
      return nil, nil, false
   }
   pre, buf := r.buf[:low], r.buf[low:high]
   r.buf = r.buf[high:]
   return pre, buf, true
}
