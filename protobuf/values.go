package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func consumeBytes(b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeBytes(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return nil, err
   }
   return val, nil
}

func consumeField(b []byte) (protowire.Number, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(b)
   err := protowire.ParseError(fLen)
   if err != nil {
      return 0, 0, 0, err
   }
   return num, typ, fLen, nil
}

func consumeFixed32(b []byte) (uint32, error) {
   val, vLen := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeFixed64(b []byte) (uint64, error) {
   val, vLen := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeGroup(num protowire.Number, b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeGroup(num, b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return nil, err
   }
   return val, nil
}

func consumeTag(b []byte) (int, error) {
   _, _, tLen := protowire.ConsumeTag(b)
   err := protowire.ParseError(tLen)
   if err != nil {
      return 0, err
   }
   return tLen, nil
}

func consumeVarint(b []byte) (uint64, error) {
   val, vLen := protowire.ConsumeVarint(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func (m Message) add(key protowire.Number, val Message) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{typ, val}
   case []Message:
      m[tag] = append(typ, val)
   }
}

func (m Message) addBytes(key protowire.Number, val []byte) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case []byte:
      m[tag] = [][]byte{typ, val}
   case [][]byte:
      m[tag] = append(typ, val)
   }
}

func (m Message) addString(key protowire.Number, val string) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{typ, val}
   case []string:
      m[tag] = append(typ, val)
   }
}

func (m Message) addUint32(key protowire.Number, val uint32) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint32:
      m[tag] = []uint32{typ, val}
   case []uint32:
      m[tag] = append(typ, val)
   }
}

func (m Message) addUint64(key protowire.Number, val uint64) {
   tag := Tag{Number: key}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{typ, val}
   case []uint64:
      m[tag] = append(typ, val)
   }
}

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

func (m Message) GetBytes(keys ...protowire.Number) []byte {
   for _, key := range keys {
      tag := Tag{Number: key}
      switch val := m[tag].(type) {
      case Message:
         m = val
      case []byte:
         return val
      }
   }
   return nil
}

func (m Message) GetMessages(keys ...protowire.Number) []Message {
   for _, key := range keys {
      tag := Tag{Number: key}
      switch val := m[tag].(type) {
      case Message:
         m = val
      case []Message:
         return val
      }
   }
   return nil
}

func (m Message) GetString(keys ...protowire.Number) string {
   for _, key := range keys {
      tag := Tag{Number: key}
      switch val := m[tag].(type) {
      case Message:
         m = val
      case string:
         return val
      }
   }
   return ""
}

func (m Message) GetUint64(keys ...protowire.Number) uint64 {
   for _, key := range keys {
      tag := Tag{Number: key}
      switch val := m[tag].(type) {
      case Message:
         m = val
      case uint64:
         return val
      }
   }
   return 0
}

func (m Message) Set(key protowire.Number, val Message) {
   tag := Tag{Number: key}
   m[tag] = val
}

func (m Message) SetUint64(key protowire.Number, val uint64) {
   tag := Tag{Number: key}
   m[tag] = val
}
