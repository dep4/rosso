package strconv

import (
   "strconv"
)

type Buffer []byte

func (b *Buffer) Byte(val byte) {
   *b = append(*b, val)
}

func (b *Buffer) Int(val int64) {
   *b = strconv.AppendInt(*b, val, 10)
}

func (b *Buffer) Quote(val string) {
   *b = strconv.AppendQuote(*b, val)
}

func (b *Buffer) Uint(val uint64) {
   *b = strconv.AppendUint(*b, val, 10)
}
