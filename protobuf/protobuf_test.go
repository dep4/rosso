package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestProto(t *testing.T) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   toks := Parse(data)
   fmt.Printf("%+v\n", toks)
}
