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

func (m Message) get(num Number) Message {
   switch value := m[num].(type) {
   case Message:
      return value
   case string:
      return m.get(-num)
   }
   return nil
}

func (m Message) value(nums ...Number) any {
   for i, num := range nums {
      if i == len(nums)-1 {
         return m[num]
      } else {
         m = m.get(num)
      }
   }
   return nil
}
