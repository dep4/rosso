package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) Add(num protowire.Number, name string, val Message) error {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{value, val}
   case []Message:
      m[tag] = append(value, val)
   }
   return nil
}

func (m Message) addString(num protowire.Number, val string) {
   tag := Tag{num, stringType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{value, val}
   case []string:
      m[tag] = append(value, val)
   }
}

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
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
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []string:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []Message:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   }
   return buf
}

func (m Message) Marshal() []byte {
   var buf []byte
   for tag, val := range m {
      buf = appendField(buf, tag.Number, val)
   }
   return buf
}
