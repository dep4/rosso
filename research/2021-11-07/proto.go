package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case Repeated:
      for _, v := range val {
         buf = appendField(buf, num, v)
      }
   case bool:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, protowire.EncodeBool(val))
   case int32:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, uint64(val))
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   }
   return buf
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, val := range m {
      buf = appendField(buf, num, val)
   }
   return buf
}
