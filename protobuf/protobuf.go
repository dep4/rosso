package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func consume(n protowire.Number, t protowire.Type, data []byte) (interface{}, int) {
   switch t {
   case protowire.VarintType:
      return protowire.ConsumeVarint(data)
   case protowire.Fixed32Type:
      return protowire.ConsumeFixed32(data)
   case protowire.Fixed64Type:
      return protowire.ConsumeFixed64(data)
   case protowire.BytesType:
      v, vLen := protowire.ConsumeBytes(data)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return string(v), vLen
   case protowire.StartGroupType:
      v, vLen := protowire.ConsumeGroup(n, data)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return v, vLen
   }
   return nil, 0
}

type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value interface{}
}

func Parse(data []byte) []Field {
   var flds []Field
   for len(data) > 0 {
      n, t, fLen := protowire.ConsumeField(data)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(data[:fLen])
      if tLen <= 0 {
         return nil
      }
      v, vLen := consume(n, t, data[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      flds = append(flds, Field{n, t, v})
      data = data[fLen:]
   }
   return flds
}
