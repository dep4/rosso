package strconv

import (
   "strconv"
   "unicode/utf8"
)

func AppendInt[T Signed](dst []byte, i T, base int) []byte {
   return strconv.AppendInt(dst, int64(i), base)
}

func AppendLabel[T Ordered](dst []byte, value T, u Unit) []byte {
   u.Factor *= float64(value)
   dst = strconv.AppendFloat(dst, u.Factor, 'f', 3, 64)
   return append(dst, u.Name...)
}

func AppendQuote[T String](dst []byte, s T) []byte {
   return strconv.AppendQuote(dst, string(s))
}

func AppendScale[T Ordered](dst []byte, value T, units []Unit) []byte {
   var u Unit
   for _, u = range units {
      if float64(value) * u.Factor < 1000 {
         break
      }
   }
   return AppendLabel(dst, value, u)
}

func AppendUint[T Unsigned](dst []byte, i T, base int) []byte {
   return strconv.AppendUint(dst, uint64(i), base)
}

func Ratio[T, U Ordered](value T, total U) float64 {
   var r float64
   if total != 0 {
      r = float64(value) / float64(total)
   }
   return r
}

// mimesniff.spec.whatwg.org#binary-data-byte
func Valid(p []byte) bool {
   for _, b := range p {
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
   return utf8.Valid(p)
}

type Ordered interface {
   Signed | Unsigned | ~float32 | ~float64
}

type Signed interface {
   ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type String interface {
   ~[]byte | ~[]rune | ~byte | ~rune | ~string
}

type Unit struct {
   Factor float64
   Name string
}

type Unsigned interface {
   ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
