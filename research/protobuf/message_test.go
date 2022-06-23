package protobuf

import (
   "fmt"
   "testing"
)

func Test_ProtoBuf(t *testing.T) {
   mes := make(message)
   fmt.Println(mes.add_varint(1, 2))
   fmt.Println(mes.add_string(1, "three"))
   fmt.Println(mes.add_varint(1, 4))
   fmt.Println(mes.add_string(1, "five"))
   fmt.Println(mes.add_varint(1, 6))
   fmt.Println(mes)
}
