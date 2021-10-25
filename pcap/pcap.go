package pcap

import (
   "bytes"
   "encoding/binary"
   "github.com/refraction-networking/utls"
)

type Handshake []byte

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
