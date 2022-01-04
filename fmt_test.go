package format

import (
   "fmt"
   "os"
   "testing"
)

func TestFmt(t *testing.T) {
   percent(os.Stdout, 2, 3)
}
