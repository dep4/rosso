package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

// Its kind of silly to include a one line function, but I keep forgetting how
// to indent the result of `Parse`.
func Indent(v interface{}) ([]byte, error) {
   return json.MarshalIndent(v, "", " ")
}

func consume(n protowire.Number, t protowire.Type, buf []byte) (interface{}, int) {
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
      v, vLen := protowire.ConsumeGroup(n, buf)
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

func Parse(buf []byte) []Field {
   var fields []Field
   for len(buf) > 0 {
      n, t, fLen := protowire.ConsumeField(buf)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if tLen <= 0 {
         return nil
      }
      v, vLen := consume(n, t, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      fields = append(fields, Field{n, t, v})
      buf = buf[fLen:]
   }
   return fields
}

