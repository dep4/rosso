package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(out []byte, key protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case array:
      for _, v := range val {
         out = appendField(out, key, v)
      }
   case bool:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, protowire.EncodeBool(val))
   case object:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendBytes(out, val.marshal())
   case string:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendString(out, val)
   case uint64:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, val)
   }
   return out
}

type array []interface{}

type object map[protowire.Number]interface{}

func (o object) marshal() []byte {
   var out []byte
   for key, val := range o {
      out = appendField(out, key, val)
   }
   return out
}
