package strings

import (
   "fmt"
   "testing"
)

func TestReference(t *testing.T) {
   var b Builder
   fmt.Fprint(&b, "hello world")
   fmt.Println(b.String())
}

func TestValue(t *testing.T) {
   var b Builder
   b.WriteByte('!')
   b.WriteRune('😀')
   b.WriteString("hello")
   b.WriteInt(9, 10)
   b.WriteQuote("world")
   b.WriteUint(9, 10)
   fmt.Println(b.String())
}
