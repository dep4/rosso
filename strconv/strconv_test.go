package strconv

import (
   "bytes"
   "fmt"
   "os"
   "strconv"
   "strings"
   "testing"
)

const (
   runs = 199_999
   value byte = 3
)

func Benchmark_Format(b *testing.B) {
   var s string
   for n := 0; n < runs; n++ {
      s += strconv.FormatInt(1, 10)
   }
}

func Benchmark_Build(b *testing.B) {
   var s strings.Builder
   for n := 0; n < runs; n++ {
      s.WriteString(strconv.FormatInt(1, 10))
   }
}

func Benchmark_Append(b *testing.B) {
   var s []byte
   for n := 0; n < runs; n++ {
      s = strconv.AppendInt(s, 1, 10)
   }
}

func Test_String(t *testing.T) {
   var str string
   str += FormatInt(2, 10)
   str += FormatUint(value, 10)
   str += Number(4444)
   str += Percent(5, 6)
   str += Quote("world")
   fmt.Println(str)
}

func Test_Writer(t *testing.T) {
   buf := new(bytes.Buffer)
   buf.WriteString(FormatInt(2, 10))
   buf.WriteString(FormatUint(value, 10))
   buf.WriteString(Number(4444))
   buf.WriteString(Percent(5, 6))
   buf.WriteString(Quote("world"))
   buf.WriteByte('\n')
   os.Stdout.ReadFrom(buf)
}
