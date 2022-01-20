package strings

import (
   "strconv"
)

type Builder []byte

func (b *Builder) Add(s string) {
   *b = append(*b, s...)
}

func (b *Builder) AddInt64(i int64, base int) {
   *b = strconv.AppendInt(*b, i, base)
}

func (b *Builder) AddQuote(s string) {
   *b = strconv.AppendQuote(*b, s)
}

func (b *Builder) AddUint64(i uint64, base int) {
   *b = strconv.AppendUint(*b, i, base)
}
