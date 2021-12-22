package crypto

import (
   "fmt"
   "testing"
)

func TestBytes(t *testing.T) {
   var b []byte
   b = append(b, 0,0,0,5, 'h', 'e', 'l', 'l', 'o')
   b = append(b, 0,0,0,5, 'w', 'o', 'r', 'l', 'd')
   buf := NewBuffer(b)
   one, two, ok := buf.ReadUint32LengthPrefixed()
   fmt.Printf("%v %s %v\n", one, two, ok)
   one, two, ok = buf.ReadUint32LengthPrefixed()
   fmt.Printf("%v %s %v\n", one, two, ok)
}
