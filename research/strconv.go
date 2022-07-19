package strconv

import (
   "strconv"
   "unicode/utf8"
)

func AppendCardinal[T Integer](dst []byte, value T) []byte {
   units := []unit{
      {1e-3, " thousand"},
      {1e-6, " million"},
      {1e-9, " billion"},
      {1e-12, " trillion"},
   }
   return scale(dst, value, units)
}

func AppendInt[T Signed](dst []byte, i T, base int) []byte {
   return strconv.AppendInt(dst, int64(i), base)
}

func AppendPercent(dst []byte, value float64) []byte {
   return label(dst, value, unit{100, "%"})
}

func AppendQuote[T String](dst []byte, s T) []byte {
   return strconv.AppendQuote(dst, string(s))
}

func AppendRate(dst []byte, value float64) []byte {
   units := []unit{
      {1e-3, " kilobyte/s"},
      {1e-6, " megabyte/s"},
      {1e-9, " gigabyte/s"},
      {1e-12, " terabyte/s"},
   }
   return scale(dst, value, units)
}

func AppendSize[T Integer](dst []byte, value T) []byte {
   units := []unit{
      {1e-3, " kilobyte"},
      {1e-6, " megabyte"},
      {1e-9, " gigabyte"},
      {1e-12, " terabyte"},
   }
   return scale(dst, value, units)
}

func AppendUint[T Unsigned](dst []byte, i T, base int) []byte {
   return strconv.AppendUint(dst, uint64(i), base)
}

func Ratio[T, U Integer](value T, total U) float64 {
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

func label[T Ordered](dst []byte, value T, u unit) []byte {
   u.factor *= float64(value)
   dst = strconv.AppendFloat(dst, u.factor, 'f', 3, 64)
   return append(dst, u.name...)
}

func scale[T Ordered](dst []byte, value T, units []unit) []byte {
   var u unit
   for _, u = range units {
      if u.factor * float64(value) < 1000 {
         break
      }
   }
   return label(dst, value, u)
}

type Integer interface {
   Signed | Unsigned
}

type Ordered interface {
   Integer | ~float32 | ~float64
}

type Signed interface {
   ~int | ~int8 | ~int16 | ~int32 | ~int64
}

type String interface {
   ~[]byte | ~[]rune | ~byte | ~rune | ~string
}

type Unsigned interface {
   ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type unit struct {
   factor float64
   name string
}
