package tls

import (
   "testing"
)

func TestMarshal(t *testing.T) {
   h, err := Parse(Android)
   if err != nil {
      t.Fatal(err)
   }
   j, err := Marshal(h)
   if err != nil {
      t.Fatal(err)
   }
   if j != Android {
      t.Fatal(j)
   }
}
