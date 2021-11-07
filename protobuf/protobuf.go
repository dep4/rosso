package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(out []byte, key protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case Array:
      for _, v := range val {
         out = appendField(out, key, v)
      }
   case bool:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, protowire.EncodeBool(val))
   case Object:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendBytes(out, val.Marshal())
   case string:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendString(out, val)
   case uint64:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, val)
   }
   return out
}

func consume(key protowire.Number, typ protowire.Type, buf []byte) (interface{}, int) {
   switch typ {
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
      v, vLen := protowire.ConsumeGroup(key, buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return v, vLen
   }
   return nil, 0
}

type Array []interface{}

type Object map[protowire.Number]interface{}

func Parse(buf []byte) Object {
   fs := make(map[protowire.Number]interface{})
   for len(buf) > 0 {
      key, typ, fLen := protowire.ConsumeField(buf)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if tLen <= 0 {
         return nil
      }
      v, vLen := consume(key, typ, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      alfa, ok := fs[key]
      if ok {
         bravo, ok := alfa.([]interface{})
         if ok {
            fs[key] = append(bravo, v)
         } else {
            fs[key] = []interface{}{alfa, v}
         }
      } else {
         fs[key] = v
      }
      buf = buf[fLen:]
   }
   return fs
}

func (o Object) Marshal() []byte {
   var out []byte
   for key, val := range o {
      out = appendField(out, key, val)
   }
   return out
}
