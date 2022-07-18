package strconv

import (
   "fmt"
   "strconv"
   "strings"
   "testing"
)

const (
   runs = 199_999
   value int64 = 9
)

func Benchmark_Strconv(b *testing.B) {
   var s string
   for n := 0; n < runs; n++ {
      s += strconv.FormatInt(value, 10)
   }
}

func Benchmark_Fmt(b *testing.B) {
   s := new(strings.Builder)
   for n := 0; n < runs; n++ {
      fmt.Fprint(s, value)
   }
}
