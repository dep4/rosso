package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

const (
   BytesType Type = 2
   Fixed64Type Type = 1
   MessageType Type = 0
   VarintType Type = 6
)

type Type = protowire.Type

func appendField(buf []byte, num Number, val any) []byte {
   switch val := val.(type) {
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []byte:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val)
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case []uint64:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []string:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []Message:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   }
   return buf
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
      value := buf[tLen:fLen]
      switch typ {
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         Add(mes, num, VarintType, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         Add(mes, num, Fixed64Type, val)
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(value)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         binary := format.IsBinary(val)
         mes, err := Unmarshal(val)
         if err != nil {
            if binary {
               Add(mes, num, BytesType, val)
            } else {
               Add(mes, num, BytesType, string(val))
            }
         } else {
            Add(mes, num, MessageType, mes)
            if !binary {
               Add(mes, num, BytesType, string(val))
            }
         }
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Marshal() []byte {
   var buf []byte
   for tag, value := range m {
      buf = appendField(buf, tag.Number, value)
   }
   return buf
}

type Number = protowire.Number

type Tag struct {
   Number
   Type
}

func NewTag(num Number) Tag {
   return Tag{Number: num}
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

type Message map[Tag]any

func Add[T any](mes Message, num Number, typ Type, val T) {
   key := Tag{num, typ}
   switch value := mes[key].(type) {
   case nil:
      mes[key] = val
   case T:
      mes[key] = []T{value, val}
   case []T:
      mes[key] = append(value, val)
   }
}

////////////////////////////////////////////////////////////////////////////////

type Token[T any] struct {
   Message
   Value T
}

func newToken[T any](m Message) Token[T] {
   return Token[T]{Message: m}
}

func (t Token[T]) get(key Tag) Token[T] {
   switch val := t.Message[key].(type) {
   case Message:
      t.Message = val
   case T:
      t.Value = val
   }
   return t
}
