package strconv

import (
   "strconv"
)

type Buffer []byte

func (b *Buffer) AppendInt(i int64) {
   *b = strconv.AppendInt(*b, i, 10)
}

func (b *Buffer) AppendQuote(val string) {
   *b = strconv.AppendQuote(*b, val)
}

func (b *Buffer) AppendUint(val uint64) {
   *b = strconv.AppendUint(*b, val, 10)
}

func (b *Buffer) WriteByte(c byte) {
   *b = append(*b, c)
}

func (b *Buffer) WriteString(s string) {
   *b = append(*b, s...)
}

func Number[T Ordered](value T) string {
   return label(value, "", " K", " M", " B", " T")
}

func Percent[T, U Ordered](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

func Rate[T, U Ordered](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = float64(value) / float64(total)
   }
   return label(ratio, " B/s", " kB/s", " MB/s", " GB/s", " TB/s")
}

func Size[T Ordered](value T) string {
   return label(value, " B", " kB", " MB", " GB", " TB")
}

func label[T Ordered](value T, units ...string) string {
   var (
      i int
      unit string
      val = float64(value)
   )
   for i, unit = range units {
      if val < 1000 {
         break
      }
      val /= 1000
   }
   if i >= 1 {
      i = 3
   }
   return strconv.FormatFloat(val, 'f', i, 64) + unit
}

type Ordered interface {
   float64 | int | int64 | uint64
}
