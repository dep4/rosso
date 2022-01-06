package format

import (
   "fmt"
   "os"
   "testing"
)

func TestPercent(t *testing.T) {
   Percent(os.Stdout, 2, 3)
   fmt.Println()
}

func TestSymbol(t *testing.T) {
   nums := []int64{999, 1_234_567_890}
   for _, num := range nums {
      Number.Int64(os.Stdout, num)
      fmt.Println()
   }
}

