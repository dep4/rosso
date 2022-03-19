package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Tag struct {
   protowire.Number
   protowire.Type
}

const (
   messageType = 0
   fixed64Type = 1
   bytesType = 2
   varintType = 6
)

type Message map[Tag]any

func (m Message) GetMessages(num protowire.Number, s string) []Message {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case []Message:
      return value
   case Message:
      return []Message{value}
   }
   return nil
}

func (m Message) GetFixed64(num protowire.Number, s string) uint64 {
   tag := Tag{num, fixed64Type}
   value, ok := m[tag].(uint64)
   if ok {
      return value
   }
   return 0
}

func (m Message) GetString(num protowire.Number, s string) string {
   tag := Tag{num, bytesType}
   value, ok := m[tag].(string)
   if ok {
      return value
   }
   return ""
}

func (m Message) GetVarint(num protowire.Number, s string) uint64 {
   tag := Tag{num, varintType}
   value, ok := m[tag].(uint64)
   if ok {
      return value
   }
   return 0
}
