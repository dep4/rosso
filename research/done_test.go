package strconv

import (
   "os"
   "testing"
)

func Test_Append(t *testing.T) {
   var dst []byte
   dst = NewRatio(2, 3).AppendPercent(dst)
   dst = append(dst, '\n')
   dst = AppendCardinal(dst, 9999)
   dst = append(dst, '\n')
   dst = NewRatio(12345, 10).AppendRate(dst)
   dst = append(dst, '\n')
   os.Stdout.Write(dst)
}
