package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

type Number = protowire.Number

type Message map[Number]any

func Unmarshal(in []byte) (Message, error) {
   mes := make(Message)
   for len(in) >= 1 {
      num, typ, fLen := protowire.ConsumeField(in)
      if err := protowire.ParseError(fLen); err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(in[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      buf := in[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         if len(val) >= 1 {
            embed, err := Unmarshal(val)
            if err != nil {
               add(mes, num, string(val))
            } else if format.IsBinary(val) {
               add(mes, num, embed)
            } else {
               add(mes, num, string(val))
               add(mes, -num, embed)
            }
         } else {
            add(mes, num, "")
         }
      case protowire.Fixed32Type:
         val, vLen := protowire.ConsumeFixed32(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, val)
      }
      in = in[fLen:]
   }
   return mes, nil
}

func add[T any](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = []T{value, val}
   case []T:
      mes[num] = append(value, val)
   }
}
