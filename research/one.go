package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

type Token interface {
   Message | []byte | string | uint64
}

func add[T Token](m Message, num Number, typ Type, val T) {
   key := Tag{num, typ}
   switch value := m[key].(type) {
   case nil:
      m[key] = val
   case T:
      m[key] = []T{value, val}
   case []T:
      m[key] = append(value, val)
   }
}

type Message map[Tag]interface{}

type Number = protowire.Number

type Tag struct {
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

func Unmarshal(buf []byte) (Message, error) {
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      err := protowire.ParseError(fLen)
      if err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      value := buf[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         if err := mes.consumeBytes(num, value); err != nil {
            return nil, err
         }
      case protowire.Fixed64Type:
         if err := mes.consumeFixed64(num, value); err != nil {
            return nil, err
         }
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, varintType, val)
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) consumeBytes(num Number, b []byte) error {
   val, vLen := protowire.ConsumeBytes(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         add(m, num, bytesType, val)
      } else {
         add(m, num, bytesType, string(val))
      }
   } else {
      add(m, num, messageType, mes)
      if !binary {
         add(m, num, bytesType, string(val))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num Number, b []byte) error {
   val, vLen := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   add(m, num, fixed64Type, val)
   return nil
}
