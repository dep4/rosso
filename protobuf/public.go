package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func Add[T any](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = []T{value, val}
   case []T:
      mes[num] = append(value, val)
   }
}

func Value[T any](mes Message, nums ...Number) T {
   val, _ := mes.value(nums...).(T)
   return val
}

func Values[T any](mes Message, nums ...Number) []T {
   switch value := mes.value(nums...).(type) {
   case []T:
      return value
   case T:
      return []T{value}
   }
   return nil
}

type Message map[Number]any

func Decode(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

func Unmarshal(buf []byte) (Message, error) {
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      if err := protowire.ParseError(fLen); err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      field := buf[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         if len(val) >= 1 {
            value, err := Unmarshal(val)
            if err != nil {
               Add(mes, num, string(val))
            } else if format.IsBinary(val) {
               Add(mes, num, value)
            } else {
               Add(mes, num, string(val))
               Add(mes, -num, value)
            }
         } else {
            Add(mes, num, "")
         }
      case protowire.Fixed32Type:
         val, vLen := protowire.ConsumeFixed32(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         Add(mes, num, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         Add(mes, num, val)
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         Add(mes, num, val)
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, val := range m {
      if num >= protowire.MinValidNumber {
         buf = appendField(buf, num, val)
      }
   }
   return buf
}

type Number = protowire.Number
