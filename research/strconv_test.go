package strconv

import (
   "bytes"
   "os"
   "testing"
)

func Test_Write(t *testing.T) {
   buf := new(bytes.Buffer)
   buf.WriteRune('😀')
   buf.WriteString("hello")
   WriteInt(buf, 9, 10)
   WriteUint(buf, 8, 10)
   WriteQuote(buf, "world")
   WritePercent(buf, 2, 3)
   WriteNumber(buf, 9999)
   buf.WriteByte('\n')
   os.Stdout.ReadFrom(buf)
}
