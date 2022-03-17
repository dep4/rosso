package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

type token interface {
   message | []byte | string | uint64
}

func add[T token](mes message, num Number, typ Type, val T) {
   key := tag{num, typ}
   switch value := mes[key].(type) {
   case nil:
      mes[key] = val
   case T:
      mes[key] = []T{value, val}
   case []T:
      mes[key] = append(value, val)
   }
}

type message map[tag]interface{}

type Number = protowire.Number

type tag struct {
   Number
   Type
}

type Type = protowire.Type

const (
   messageType Type = 0
   fixed64Type Type = 1
   bytesType Type = 2
   varintType Type = 6
)

func unmarshal(buf []byte) (message, error) {
   mes := make(message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      if err := protowire.ParseError(fLen); err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      value := buf[tLen:fLen]
      switch typ {
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, varintType, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, fixed64Type, val)
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         binary := format.IsBinary(val)
         mes, err := unmarshal(val)
         if err != nil {
            if binary {
               add(mes, num, bytesType, val)
            } else {
               add(mes, num, bytesType, string(val))
            }
         } else {
            add(mes, num, messageType, mes)
            if !binary {
               add(mes, num, bytesType, string(val))
            }
         }
      }
      buf = buf[fLen:]
   }
   return mes, nil
}
