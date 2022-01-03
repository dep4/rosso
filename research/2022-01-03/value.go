package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Tag struct {
   protowire.Number
   String string
}

type Message map[Tag]interface{}

func (m Message) Get(n protowire.Number, s string) Message {
   val, ok := m[Tag{n, s}].(Message)
   if ok {
      return val
   }
   if val, ok := m[Tag{n, ""}].(Message); ok {
      return val
   }
   return m
}
