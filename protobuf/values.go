package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

// bytes or group
func (m Message) Add(num protowire.Number, key string, mes Message) {
   tag := Tag{num, key}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = mes
   case Message:
      m[tag] = []Message{val, mes}
   case []Message:
      m[tag] = append(val, mes)
   }
}

func (m Message) Get(num protowire.Number, key string) Message {
   for _, str := range []string{"bytes", "group"} {
      val, ok := m[Tag{num, str}].(Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number, key string) []byte {
   val, ok := m[Tag{num, "bytes"}].([]byte)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetMessages(num protowire.Number, key string) []Message {
   for _, str := range []string{"bytes", "group"} {
      val, ok := m[Tag{num, str}].([]Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetString(num protowire.Number, key string) string {
   val, ok := m[Tag{num, "bytes"}].(string)
   if ok {
      return val
   }
   return ""
}

func (m Message) GetUint64(num protowire.Number, key string) uint64 {
   for _, str := range []string{"varint", "fixed64"} {
      val, ok := m[Tag{num, str}].(uint64)
      if ok {
         return val
      }
   }
   return 0
}

// bytes
func (m Message) addBytes(num protowire.Number, key string, v []byte) {
   tag := Tag{num, key}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case []byte:
      m[tag] = [][]byte{val, v}
   case [][]byte:
      m[tag] = append(val, v)
   }
}

// bytes
func (m Message) addString(num protowire.Number, key, v string) {
   tag := Tag{num, key}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case string:
      m[tag] = []string{val, v}
   case []string:
      m[tag] = append(val, v)
   }
}

// fixed32
func (m Message) addUint32(num protowire.Number, key string, v uint32) {
   tag := Tag{num, key}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint32:
      m[tag] = []uint32{val, v}
   case []uint32:
      m[tag] = append(val, v)
   }
}

// varint or fixed64
func (m Message) addUint64(num protowire.Number, key string, v uint64) {
   tag := Tag{num, key}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint64:
      m[tag] = []uint64{val, v}
   case []uint64:
      m[tag] = append(val, v)
   }
}

func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   ok := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if ok {
         m.addBytes(num, "bytes", val)
      } else {
         m.addString(num, "bytes", string(val))
      }
   } else if ok {
      // Could be Message or []byte
      m.Add(num, "bytes", mes)
   } else {
      // Cound be Message or string
      m.addString(num, "bytes", string(val))
   }
   return nil
}

func (m Message) consumeFixed32(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint32(num, "fixed32", val)
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, "fixed64", val)
   return nil
}

func (m Message) consumeGroup(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeGroup(num, buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   mes, err := Unmarshal(val)
   if err != nil {
      return err
   }
   m.Add(num, "group", mes)
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, "varint", val)
   return nil
}
