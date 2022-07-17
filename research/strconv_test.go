package strconv

import (
   "fmt"
   "io"
   "strconv"
   "strings"
   "testing"
)

var value int64 = 9

func Benchmark_Fmt(b *testing.B) {
   s := new(strings.Builder)
   for n := 0; n < b.N; n++ {
      fmt.Fprint(s, value)
   }
}

func Benchmark_Rosso(b *testing.B) {
   s := new(strings.Builder)
   for n := 0; n < b.N; n++ {
      fprint(s, value)
   }
}

type signed interface {
   ~int | ~int8 | ~int16 | ~int32 | ~int64
}

func fprint[T signed](w io.Writer, value T) (int, error) {
   return io.WriteString(w, strconv.FormatInt(int64(value), 10))
}
