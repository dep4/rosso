package strconv

import (
   "os"
   "testing"
)

func Test_Append(t *testing.T) {
   var b []byte
   b = AppendCardinal(b, 1234)
   b = append(b, '\n')
   b = AppendInt(b, 9, 10)
   b = append(b, '\n')
   b = AppendQuote(b, "hello")
   b = append(b, '\n')
   b = AppendSize(b, 1234)
   b = append(b, '\n')
   b = AppendUint[byte](b, 9, 10)
   b = append(b, '\n')
   b = NewRatio(12345, 10).AppendCardinal(b)
   b = append(b, '\n')
   b = NewRatio(2, 3).AppendPercent(b)
   b = append(b, '\n')
   b = NewRatio(12345, 10).AppendRate(b)
   b = append(b, '\n')
   os.Stdout.Write(b)
}
