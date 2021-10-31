package bytes

import (
   "bytes"
   "encoding/binary"
)

type buffer struct {
   buf []byte
}

// godocs.io/bytes#Buffer.Next
func (b *buffer) next(n int) ([]byte, bool) {
   if n < 0 || n > len(b.buf) {
      return nil, false
   }
   buf := b.buf[:n]
   b.buf = b.buf[n:]
   return buf, true
}

// godocs.io/bytes#Buffer.ReadBytes
func (b *buffer) readBytes(delim byte) ([]byte, bool) {
   cut := bytes.IndexByte(b.buf, delim) + 1
   if cut == 0 {
      return nil, false
   }
   buf := b.buf[:cut]
   b.buf = b.buf[cut:]
   return buf, true
}

// godocs.io/golang.org/x/crypto/cryptobyte#String.ReadUint16LengthPrefixed
func (b *buffer) readUint16LengthPrefixed() ([]byte, []byte, bool) {
   if len(b.buf) < 2 {
      return nil, nil, false
   }
   high := 2 + binary.BigEndian.Uint16(b.buf)
   if len(b.buf) < int(high) {
      return nil, nil, false
   }
   pre, buf := b.buf[:2], b.buf[2:high]
   b.buf = b.buf[high:]
   return pre, buf, true
}
