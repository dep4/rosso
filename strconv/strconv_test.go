package strconv

import (
   "bytes"
   "fmt"
   "os"
   "testing"
)

var value byte = 3

func Test_String(t *testing.T) {
   var buf string
   buf += FormatInt(2, 10)
   buf += FormatUint(value, 10)
   buf += Number(4444)
   buf += Percent(5, 6)
   buf += Quote("world")
   fmt.Println(buf)
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
