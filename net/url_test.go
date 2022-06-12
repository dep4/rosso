package net

import (
   "fmt"
   "testing"
)

func TestValues(t *testing.T) {
   text := []byte("alfa=bravo&charlie=delta")
   var val Values
   err := val.UnmarshalText(text)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(val.Get("charlie"))
}
