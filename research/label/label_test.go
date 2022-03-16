package format

import (
   "fmt"
   "testing"
)

func TestMeasure(t *testing.T) {
   label := new(measure[int]).labelNumber(9_999)
   fmt.Println(label)
}
