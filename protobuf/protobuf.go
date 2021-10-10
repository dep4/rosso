package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value interface{}
}

func ParseUnknown(data []byte) []Field {
   var flds []Field
   for len(data) > 0 {
      n, t, fLen := protowire.ConsumeField(data)
      if fLen < 1 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(data[:fLen])
      if tLen < 1 {
         return nil
      }
      var (
         v interface{}
         vLen int
      )
      switch t {
      case protowire.VarintType:
         v, vLen = protowire.ConsumeVarint(data[tLen:fLen])
      case protowire.Fixed32Type:
         v, vLen = protowire.ConsumeFixed32(data[tLen:fLen])
      case protowire.Fixed64Type:
         v, vLen = protowire.ConsumeFixed64(data[tLen:fLen])
      case protowire.BytesType:
         v, vLen = protowire.ConsumeBytes(data[tLen:fLen])
         sub := ParseUnknown(v.([]byte))
         if sub != nil {
            v = sub
         }
      case protowire.StartGroupType:
         v, vLen = protowire.ConsumeGroup(n, data[tLen:fLen])
         sub := ParseUnknown(v.([]byte))
         if sub != nil {
            v = sub
         }
      }
      if vLen < 1 {
         return nil
      }
      fld := Field{Number: n, Type: t}
      bytes, ok := v.([]byte)
      if ok {
         fld.Value = string(bytes)
      } else {
         fld.Value = v
      }
      flds = append(flds, fld)
      data = data[fLen:]
   }
   return flds
}
