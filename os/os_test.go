package os

import (
   "github.com/89z/rosso/strconv"
   "os"
   "testing"
)

func Test_Create(t *testing.T) {
   file, err := Create("ignore.txt")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
}

func Test_Progress(t *testing.T) {
   var b []byte
   b = strconv.NewRatio(1234, 10000).AppendPercent(b)
   b = append(b, "   "...)
   b = strconv.AppendSize(b, 1234)
   b = append(b, "   "...)
   b = strconv.NewRatio(123456, 100).AppendRate(b)
   b = append(b, '\n')
   os.Stderr.Write(b)
}
