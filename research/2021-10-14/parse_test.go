package parse

import (
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

var BytesType = protopack.BytesType

type (
   Bytes = protopack.Bytes
   LengthPrefix = protopack.LengthPrefix
   Message = protopack.Message
   Tag = protopack.Tag
)

func TestParse(t *testing.T) {
   in := Message{
      Tag{6, BytesType}, LengthPrefix{
         Tag{1, BytesType}, Bytes("Vary"),
      },
      Tag{6, BytesType}, LengthPrefix{
         Tag{1, BytesType}, Bytes("MD-Version"),
      },
   }
   data := in.Marshal()
   flds := parse(data)
   fmt.Printf("%+v\n", flds)
}
