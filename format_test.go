package format

import (
   "fmt"
   "testing"
)

func TestPercent(t *testing.T) {
   tots := []int64{0, 3}
   for _, tot := range tots {
      val := Percent(2, tot)
      fmt.Println(val)
   }
}

func TestSymbol(t *testing.T) {
   nums := []int64{999, 1_234_567_890}
   for _, num := range nums {
      val := Number.LabelInt(num)
      fmt.Println(val)
   }
}
