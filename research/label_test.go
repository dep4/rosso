package strconv

import (
   "fmt"
   "testing"
)

var tests = []float64{999, 9999}

func Test_Factor(t *testing.T) {
   for _, test := range tests {
      str := Label(test, Cardinal)
      fmt.Println(str)
   }
   str := Ratio(12345, 10, Rate)
   fmt.Println(str)
}
