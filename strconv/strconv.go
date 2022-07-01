package strconv

import (
   "strconv"
   "unicode/utf8"
)

var (
   AppendQuote = strconv.AppendQuote
   AppendUint = strconv.AppendUint
   FormatFloat = strconv.FormatFloat
   Quote = strconv.Quote
)

func Number[T Ordered](value T) string {
   return label(value, "", " K", " M", " B", " T")
}

func Size[T Ordered](value T) string {
   return label(value, " B", " kB", " MB", " GB", " TB")
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

func label[T Ordered](value T, unit ...string) string {
   var (
      i int
      symbol string
      val = float64(value)
   )
   for i, symbol = range unit {
      if val < 1000 {
         break
      }
      val /= 1000
   }
   if i >= 1 {
      i = 3
   }
   return strconv.FormatFloat(val, 'f', i, 64) + symbol
}

func percent[T, U Ordered](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

func rate[T, U Ordered](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = float64(value) / float64(total)
   }
   return label(ratio, " B/s", " kB/s", " MB/s", " GB/s", " TB/s")
}

type Ordered interface {
   float64 | int | int64 | uint64
}
