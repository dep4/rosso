package bytes

import (
   "bytes"
   "encoding/binary"
   "github.com/89z/parse/tls"
)

func handshake(b []byte) *tls.ClientHello {
   for {
      low := bytes.IndexByte(b, 0x16)
      if low == -1 {
         return nil
      }
      r := newReader(b[low:])
      var w []byte
      buf, ok := r.readN(1)
      if ok {
         w = append(w, buf...)
      }
      buf, ok = r.readN(2)
      if ok {
         w = append(w, buf...)
      }
      pre, buf, ok := r.readUint16LengthPrefixed()
      if ok {
         w = append(w, pre...)
         w = append(w, buf...)
      }
      hand, err := tls.ParseHandshake(w)
      if err == nil {
         return hand
      }
      b = b[low+1:]
   }
}

type reader struct {
   b []byte
}

func newReader(b []byte) reader {
   return reader{b}
}

func (r *reader) readN(n int) ([]byte, bool) {
   if n < 0 || n > len(r.b) {
      return nil, false
   }
   buf := r.b[:n]
   r.b = r.b[n:]
   return buf, true
}

func (r *reader) readUint16LengthPrefixed() ([]byte, []byte, bool) {
   if len(r.b) < 2 {
      return nil, nil, false
   }
   high := 2 + binary.BigEndian.Uint16(r.b)
   if len(r.b) < int(high) {
      return nil, nil, false
   }
   pre, buf := r.b[:2], r.b[2:high]
   r.b = r.b[high:]
   return pre, buf, true
}
