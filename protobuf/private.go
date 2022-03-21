package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func add[T any](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = []T{value, val}
   case []T:
      mes[num] = append(value, val)
   }
}

func appendField(in []byte, num Number, val any) []byte {
   switch val := val.(type) {
   case Message:
      in = protowire.AppendTag(in, num, protowire.BytesType)
      in = protowire.AppendBytes(in, val.Marshal())
   case string:
      in = protowire.AppendTag(in, num, protowire.BytesType)
      in = protowire.AppendString(in, val)
   case uint32:
      in = protowire.AppendTag(in, num, protowire.Fixed32Type)
      in = protowire.AppendFixed32(in, val)
   case uint64:
      in = protowire.AppendTag(in, num, protowire.VarintType)
      in = protowire.AppendVarint(in, val)
   case []Message:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   case []string:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   case []uint32:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   case []uint64:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   }
   return in
}

func get[T any](mes Message, num Number) T {
   value, _ := mes[num].(T)
   return value
}
