package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func Parse(buf []byte) map[protowire.Number]interface{} {
   fs := make(map[protowire.Number]interface{})
   for len(buf) > 0 {
      k, t, fLen := protowire.ConsumeField(buf)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if tLen <= 0 {
         return nil
      }
      v, vLen := consume(k, t, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      alfa, ok := fs[k]
      if ok {
         bravo, ok := alfa.([]interface{})
         if ok {
            fs[k] = append(bravo, v)
         } else {
            fs[k] = []interface{}{alfa, v}
         }
      } else {
         fs[k] = v
      }
      buf = buf[fLen:]
   }
   return fs
}

func consume(k protowire.Number, t protowire.Type, buf []byte) (interface{}, int) {
   switch t {
   case protowire.VarintType:
      return protowire.ConsumeVarint(buf)
   case protowire.Fixed32Type:
      return protowire.ConsumeFixed32(buf)
   case protowire.Fixed64Type:
      return protowire.ConsumeFixed64(buf)
   case protowire.BytesType:
      v, vLen := protowire.ConsumeBytes(buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return string(v), vLen
   case protowire.StartGroupType:
      v, vLen := protowire.ConsumeGroup(k, buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return v, vLen
   }
   return nil, 0
}
