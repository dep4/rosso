package strconv

import (
   "os"
   "testing"
)

func Test_Label(t *testing.T) {
   var dst []byte
   dst = AppendLabel(dst, Ratio(2, 3), Percent)
   dst = append(dst, '\n')
   os.Stdout.Write(dst)
}

func Test_Scale(t *testing.T) {
   var dst []byte
   dst = AppendScale(dst, 9999, Cardinals)
   dst = append(dst, '\n')
   dst = AppendScale(dst, Ratio(12345, 10), Rates)
   dst = append(dst, '\n')
   os.Stdout.Write(dst)
}
