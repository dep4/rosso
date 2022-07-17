package strconv

import (
   "io"
   "strconv"
)

type Signed interface {
   ~int | ~int8 | ~int16 | ~int32 | ~int64
}

func FormatInt[T Signed](value T, base int) string {
   return strconv.FormatInt(int64(value), base)
}

func Int[T Signed](w io.Writer, value T, base int) (int, error) {
   return io.WriteString(w, FormatInt(value, base))
}

type Unsigned interface {
   ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func FormatUint[T Unsigned](value T, base int) string {
   return strconv.FormatUint(uint64(value), base)
}

func Uint[T Unsigned](w io.Writer, value T, base int) (int, error) {
   return io.WriteString(w, FormatUint(value, base))
}
