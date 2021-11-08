package protobuf

import (
   "fmt"
   "testing"
)

type deviceConfigurationProto struct {
   TouchScreen int32 `json:"1"`
   Keyboard int32 `json:"2"`
}

func TestEncode(t *testing.T) {
   dev := deviceConfigurationProto{
      TouchScreen: 1,
      Keyboard: 1,
   }
   enc, err := newEncoder(dev)
   if err != nil {
      t.Fatal(err)
   }
   buf, err := encode(enc)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", buf)
}
