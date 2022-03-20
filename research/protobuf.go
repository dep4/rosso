package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Number]any

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
         NewToken[uint64](mes).Add(num, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         NewToken[uint64](mes).Add(num, val)
      case protowire.Fixed32Type:
         val, vLen := protowire.ConsumeFixed32(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         NewToken[uint32](mes).Add(num, val)
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(field)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         if len(val) >= 1 {
            value, err := Unmarshal(val)
            if err != nil {
               NewToken[string](mes).Add(num, string(val))
            } else if format.IsBinary(val) {
               NewToken[Message](mes).Add(num, value)
            } else {
               NewToken[string](mes).Add(num, string(val))
               NewToken[Message](mes).Add(-num, value)
            }
         } else {
            NewToken[string](mes).Add(num, "")
         }
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

type Number = protowire.Number

type Token[T any] struct {
   Message
}

func NewToken[T any](m Message) *Token[T] {
   return &Token[T]{m}
}

func (t *Token[T]) Add(num Number, val T) {
   switch value := t.Message[num].(type) {
   case nil:
      t.Message[num] = val
   case T:
      t.Message[num] = []T{value, val}
   case []T:
      t.Message[num] = append(value, val)
   }
}

func (t Token[T]) Get(num Number) Token[T] {
   t.Message, _ = t.Message[num].(Message)
   return t
}

func (t Token[T]) Value(num Number) T {
   value, _ := t.Message[num].(T)
   return value
}

func (t Token[T]) Values(num Number) []T {
   switch value := t.Message[num].(type) {
   case []T:
      return value
   case T:
      return []T{value}
   }
   return nil
}
