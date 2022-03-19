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

type Token[T any] struct {
   Message
   protowire.Type
}

const (
   messageType = iota
   bytesType
   fixed64Type
   varintType
)
