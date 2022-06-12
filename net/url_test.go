package net

import (
   "os"
   "strings"
   "testing"
)

var r = strings.NewReader(`alfa=bravo
charlie=delta`)

func TestValues(t *testing.T) {
   val := NewValues()
   n, err := val.ReadFrom(r)
   if err != nil {
      t.Fatal(err)
   }
   if n != 24 {
      t.Fatal(n)
   }
   val.WriteTo(os.Stdout)
}
