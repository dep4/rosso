package format

import (
   "fmt"
   "testing"
)

func TestMeasure(t *testing.T) {
   fmt.Println(MeasureNumber(9_999))
   fmt.Println(MeasureSize(9_999))
   fmt.Println(MeasureRate(9_999))
}
