package strings

import (
   "strconv"
   "strings"
)

type Builder struct {
   strings.Builder
}

// func AppendInt(dst []byte, i int64, base int) []byte
func (b *Builder) WriteInt(i int64, base int) (int, error) {
   s := strconv.FormatInt(i, base)
   return b.WriteString(s)
}

// func AppendQuote(dst []byte, s string) []byte
func (b *Builder) WriteQuote(s string) (int, error) {
   s = strconv.Quote(s)
   return b.WriteString(s)
}

// func AppendUint(dst []byte, i uint64, base int) []byte
func (b *Builder) WriteUint(i uint64, base int) (int, error) {
   s := strconv.FormatUint(i, base)
   return b.WriteString(s)
}

// func AppendQuoteRune(dst []byte, r rune) []byte
