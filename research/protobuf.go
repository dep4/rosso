package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Number]any

func appendField(in []byte, num Number, val any) []byte {
   switch val := val.(type) {
   case uint32:
      in = protowire.AppendTag(in, num, protowire.Fixed32Type)
      in = protowire.AppendFixed32(in, val)
   case uint64:
      in = protowire.AppendTag(in, num, protowire.VarintType)
      in = protowire.AppendVarint(in, val)
   case string:
      in = protowire.AppendTag(in, num, protowire.BytesType)
      in = protowire.AppendString(in, val)
   case Message:
      in = protowire.AppendTag(in, num, protowire.BytesType)
      in = protowire.AppendBytes(in, val.Marshal())
   case []uint32:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   case []uint64:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   case []string:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   case []Message:
      for _, value := range val {
         in = appendField(in, num, value)
      }
   }
   return in
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, value := range m {
      if num >= protowire.MinValidNumber {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}

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

type Number = protowire.Number

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
