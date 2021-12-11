package crypto

import (
   "bytes"
   "encoding/binary"
)

type Buffer struct {
   buf []byte
}

func NewBuffer(buf []byte) *Buffer {
   return &Buffer{buf}
}

func (b *Buffer) Next(n int) ([]byte, bool) {
   if n < 0 || n > len(b.buf) {
      return nil, false
   }
   buf := b.buf[:n]
   b.buf = b.buf[n:]
   return buf, true
}

func (b *Buffer) ReadBytes(delim byte) ([]byte, bool) {
   ind := bytes.IndexByte(b.buf, delim)
   if ind == -1 {
      return nil, false
   }
   buf := b.buf[:ind+1]
   b.buf = b.buf[ind+1:]
   return buf, true
}

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