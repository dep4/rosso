package strconv

import (
   "strconv"
   "unicode/utf8"
)

func FormatInt[T Signed](value T, base int) string {
   return strconv.FormatInt(int64(value), base)
}

func FormatUint[T Unsigned](value T, base int) string {
   return strconv.FormatUint(uint64(value), base)
}

func Percent[T, U Signed](value T, total U) string {
   var ratio float64
   if total != 0 {
      ratio = 100 * float64(value) / float64(total)
   }
   return strconv.FormatFloat(ratio, 'f', 1, 64) + "%"
}

func Quote[T String](value T) string {
   return strconv.Quote(string(value))
}

// mimesniff.spec.whatwg.org#binary-data-byte
func Valid(buf []byte) bool {
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
   Signed | Unsigned | ~float32 | ~float64
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
