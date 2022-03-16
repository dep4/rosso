package format

import (
   "fmt"
   "testing"
)

func TestFormat(t *testing.T) {
   {
      s := AlfaLabel(9, "Alfa")
      fmt.Println(s)
   }
   {
      s := Bravo[int]{9}.Label("Bravo")
      fmt.Println(s)
   }
   {
      s := Charlie[int]{"Charlie"}.Label(9)
      fmt.Println(s)
   }
}
