package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

type Message map[Tag]any

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
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         NewToken[uint64](mes).Add(num, VarintType, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         NewToken[uint64](mes).Add(num, Fixed64Type, val)
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         binary := format.IsBinary(val)
         value, err := Unmarshal(val)
         if err != nil {
            if binary {
               NewToken[[]byte](mes).Add(num, BytesType, val)
            } else {
               NewToken[string](mes).Add(num, BytesType, string(val))
            }
         } else {
            NewToken[Message](mes).Add(num, MessageType, value)
            if !binary {
               NewToken[string](mes).Add(num, BytesType, string(val))
            }
         }
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

type Number = protowire.Number

type Tag struct {
   Number
   Type
}

func (t Tag) MarshalText() ([]byte, error) {
   buf := strconv.AppendInt(nil, int64(t.Number), 10)
   switch t.Type {
   case BytesType:
      buf = append(buf, " bytes"...)
   case Fixed64Type:
      buf = append(buf, " fixed64"...)
   case MessageType:
      buf = append(buf, " message"...)
   case VarintType:
      buf = append(buf, " varint"...)
   }
   return buf, nil
}

type Token[T any] struct {
   Message
   Value T
}

func NewToken[T any](mes Message) *Token[T] {
   return &Token[T]{Message: mes}
}

func (t *Token[T]) Add(num Number, typ Type, val T) {
   key := Tag{num, typ}
   switch value := t.Message[key].(type) {
   case nil:
      t.Message[key] = val
   case T:
      t.Message[key] = []T{value, val}
   case []T:
      t.Message[key] = append(value, val)
   }
}

func (t Token[T]) Get(num Number, typ Type) Token[T] {
   key := Tag{num, typ}
   switch value := t.Message[key].(type) {
   case Message:
      t.Message = value
   case T:
      t.Value = value
   }
   return t
}

type Type = protowire.Type

const (
   BytesType Type = 2
   Fixed64Type Type = 1
   MessageType Type = 0
   VarintType Type = 6
)
