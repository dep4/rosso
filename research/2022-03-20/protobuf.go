package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Tag]any

type Tag struct {
   protowire.Number
   protowire.Type
}

func (t Token[T]) Get(num protowire.Number) Token[T] {
   t.Message, _ = t.Message[Tag{num, messageType}].(Message)
   return t
}

func (t Token[T]) Value(num protowire.Number) T {
   val, _ := t.Message[Tag{num, t.Type}].(T)
   return val
}

func (t Token[T]) Values(num protowire.Number) []T {
   switch value := t.Message[Tag{num, t.Type}].(type) {
   case []T:
      return value
   case T:
      return []T{value}
   }
   return nil
}

const (
   messageType = iota
   bytesType
   fixed64Type
   varintType
)

type Token[T any] struct {
   Message
   protowire.Type
}

func Fixed64(m Message) Token[uint64] {
   return Token[uint64]{m, fixed64Type}
}

func String(m Message) Token[string] {
   return Token[string]{m, bytesType}
}

func Varint(m Message) Token[uint64] {
   return Token[uint64]{m, varintType}
}
