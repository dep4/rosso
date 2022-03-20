package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Number]any

type Number = protowire.Number

type Token[T any] struct {
   Message
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
