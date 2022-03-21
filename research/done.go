package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(buf []byte, num Number, val any) []byte {
   switch val := val.(type) {
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case uint32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case []Message:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []string:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []uint32:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []uint64:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}
