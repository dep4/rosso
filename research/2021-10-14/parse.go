package parse

import (
   "github.com/segmentio/encoding/proto"
)

type field struct {
   num proto.FieldNumber
   typ proto.WireType
   val interface{}
}

func consume(f proto.FieldNumber, t proto.WireType, dat proto.RawValue) interface{} {
   switch t {
   case proto.Fixed32:
      return dat.Fixed32()
   case proto.Fixed64:
      return dat.Fixed64()
   case proto.Varint:
      return dat.Varint()
   case proto.Varlen:
      sub := parse(dat)
      if sub != nil {
         return sub
      }
      return string(dat)
   }
   return nil
}

func parse(data []byte) []field {
   var flds []field
   for len(data) > 0 {
      f, t, dat, m, err := proto.Parse(data)
      if err != nil {
         return nil
      }
      v := consume(f, t, dat)
      if v == nil {
         return nil
      }
      flds = append(flds, field{f, t, v})
      data = m
   }
   return flds
}
