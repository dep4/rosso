package strconv

import (
   "strconv"
)

func FormatInt[T Signed](value T, base int) string {
   return strconv.FormatInt(int64(value), base)
}

type Signed interface {
   ~int | ~int8 | ~int16 | ~int32 | ~int64
}
