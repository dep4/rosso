package format

import (
   "fmt"
   "testing"
)

func TestMeasure(t *testing.T) {
   fmt.Println(LabelNumber(9_999))
   fmt.Println(LabelSize(9_999))
   fmt.Println(LabelRate(9_999))
}
