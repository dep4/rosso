package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

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

type Token[T any] struct {
   Message
   Value T
}

type Message map[Tag]any

type Tag struct {
   Number
   Type
}

type Number = protowire.Number

type Type = protowire.Type
