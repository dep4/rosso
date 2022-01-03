package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Tag]interface{}

func (m Message) Get(n protowire.Number, s string) Message {
   val, ok := m[Tag{n, s}].(Message)
   if ok {
      return val
   }
   return m
}

func (m Message) Get2(tags ...Tag) Message {
   for _, tag := range tags {
      val, ok := m[tag].(Message)
      if ok {
         m = val
      }
   }
   return m
}

func (m Message) Get3(tags []Tag) Message {
   for _, tag := range tags {
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
