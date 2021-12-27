package strconv

import (
   "fmt"
   "testing"
)

func TestSymbol(t *testing.T) {
   num := Number.FormatInt(1_234_567_890)
   fmt.Println(num)
}
