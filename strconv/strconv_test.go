package strconv

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
   num := Number.FormatInt(1_234_567_890)
   fmt.Println(num)
}
