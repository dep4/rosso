package strings

import (
   "fmt"
   "testing"
)

func TestValue(t *testing.T) {
   var b Builder
   b.Add("hello")
   b.AddInt64(9, 10)
   b.AddQuote("world")
   b.AddUint64(9, 10)
   fmt.Println(string(b))
}
