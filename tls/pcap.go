package tls

import (
   "bytes"
   "encoding/binary"
   "encoding/hex"
   "github.com/refraction-networking/utls"
   "net"
   "net/http"
   "strings"
)

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
         uconn := tls.UClient(dialConn, config, tls.HelloCustom)
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

type Handshake []byte

func Decode(s string) (Handshake, error) {
   return hex.DecodeString(s)
}

func Handshakes(data []byte) []Handshake {
   var hands []Handshake
   for {
      // start of record
      rec1 := bytes.IndexByte(data, 0x16)
      if rec1 == -1 {
         return hands
      }
      // start of version
      ver1 := rec1 + 1
      // start of length
      len1 := ver1 + 2
      // end of length
      len2 := len1 + 2
      if len2 < len(data) {
         recLen := binary.BigEndian.Uint16(data[len1:len2])
         // end of record
         rec2 := len2 + int(recLen)
         if rec2 < len(data) {
            hands = append(hands, data[rec1:rec2])
         }
      }
      data = data[rec1+1:]
   }
}

func (h Handshake) ClientHello() (*tls.ClientHelloSpec, error) {
   var fp tls.Fingerprinter
   return fp.FingerprintClientHello(h)
}
