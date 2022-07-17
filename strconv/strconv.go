package strconv

import (
   "strconv"
   "unicode/utf8"
)

type Buffer []byte

func (b *Buffer) AppendInt(i int64) {
   *b = strconv.AppendInt(*b, i, 10)
}

func (b *Buffer) AppendUint(val uint64) {
   *b = strconv.AppendUint(*b, val, 10)
}

func (b *Buffer) Write(p []byte) (int, error) {
   *b = append(*b, p...)
   return len(p), nil
}

func (b *Buffer) WriteByte(c byte) {
   *b = append(*b, c)
}

func (b *Buffer) WriteString(s string) {
   *b = append(*b, s...)
}

var FormatUint = strconv.FormatUint

func (b *Buffer) AppendQuote(val string) {
   *b = strconv.AppendQuote(*b, val)
}

// mimesniff.spec.whatwg.org#binary-data-byte
func String(buf []byte) bool {
   for _, b := range buf {
      if b <= 0x08 {
         return false
      }
      if b == 0x0B {
         return false
      }
      if b >= 0x0E && b <= 0x1A {
         return false
      }
      if b >= 0x1C && b <= 0x1F {
         return false
      }
   }
   return utf8.Valid(buf)
}

type Ordered interface {
   float64 | int | int64 | uint64
}

func Percent[T, U Ordered](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
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

func Number[T Ordered](value T) string {
   return label(value, "", " K", " M", " B", " T")
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
