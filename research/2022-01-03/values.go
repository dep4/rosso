package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Tag]interface{}

func (m Message) Get(keys ...protowire.Number) Message {
   for _, key := range keys {
      tag := Tag{Number: key}
      val, ok := m[tag].(Message)
      if ok {
         m = val
      }
   }
   return m
}

type Tag struct {
   protowire.Number
   String string
}
