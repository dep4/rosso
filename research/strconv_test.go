package strconv

import (
   "os"
   "testing"
)

func Test_Append(t *testing.T) {
   var dst []byte
   dst = AppendPercent(dst, Ratio(2, 3))
   dst = append(dst, '\n')
   dst = AppendCardinal(dst, 9999)
   dst = append(dst, '\n')
   dst = AppendRate(dst, Ratio(12345, 10))
   dst = append(dst, '\n')
   os.Stdout.Write(dst)
}
