package protobuf

import (
   "fmt"
   "testing"
)

func Test_ProtoBuf(t *testing.T) {
   mes := make(Message)
   fmt.Println(mes.Add_Varint(1, 2))
   fmt.Println(mes.Add_String(1, "three"))
   fmt.Println(mes.Add_Varint(1, 4))
   fmt.Println(mes.Add_String(1, "five"))
   fmt.Println(mes.Add_Varint(1, 6))
   fmt.Println(mes)
}
