package parse

import (
   "fmt"
   "os"
   "testing"
)

func TestParse(t *testing.T) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   flds := parse(data)
   fmt.Printf("%+v\n", flds)
}
