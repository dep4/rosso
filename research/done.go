package protobuf

import (
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

type Token[T any] struct {
   Message
   Value T
}

func NewToken[T any](mes Message) *Token[T] {
   return &Token[T]{Message: mes}
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
