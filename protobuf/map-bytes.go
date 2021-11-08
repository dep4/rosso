package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) Marshal() []byte {
   var out []byte
   for key, val := range m {
      out = appendField(out, key, val)
   }
   return out
}

func appendField(out []byte, key protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case Message:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendBytes(out, val.Marshal())
   case Repeated:
      for _, v := range val {
         out = appendField(out, key, v)
      }
   case bool:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, protowire.EncodeBool(val))
   case int32:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, uint64(val))
   case string:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendString(out, val)
   }
   return out
}
