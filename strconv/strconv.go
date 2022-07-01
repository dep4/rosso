package strconv

import (
   "strconv"
   "unicode/utf8"
)

var (
   AppendQuote = strconv.AppendQuote
   FormatFloat = strconv.FormatFloat
   Quote = strconv.Quote
)

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

type Number interface {
   float64 | int | int64 | ~uint64
}

func Label[T Number](value T, unit ...string) string {
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

func Label_Number[T Number](value T) string {
   return Label(value, "", " K", " M", " B", " T")
}

func Label_Rate[T Number](value T) string {
   return Label(value, " B/s", " kB/s", " MB/s", " GB/s", " TB/s")
}

func Label_Size[T Number](value T) string {
   return Label(value, " B", " kB", " MB", " GB", " TB")
}
