package strconv

import (
   "strconv"
)

// godocs.io/bytes#Buffer
type Buffer []byte

// godocs.io/strconv#AppendInt
func (b *Buffer) AppendInt(i int64) {
   *b = strconv.AppendInt(*b, i, 10)
}

// godocs.io/strconv#AppendQuote
func (b *Buffer) AppendQuote(val string) {
   *b = strconv.AppendQuote(*b, val)
}

// godocs.io/strconv#AppendUint
func (b *Buffer) AppendUint(val uint64) {
   *b = strconv.AppendUint(*b, val, 10)
}

// godocs.io/bytes#Buffer.Write
func (b *Buffer) Write(p []byte) (int, error) {
   *b = append(*b, p...)
   return len(p), nil
}

// godocs.io/bytes#Buffer.WriteByte
func (b *Buffer) WriteByte(c byte) {
   *b = append(*b, c)
}

// godocs.io/bytes#Buffer.WriteString
func (b *Buffer) WriteString(s string) {
   *b = append(*b, s...)
}
