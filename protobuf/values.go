package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) Get(num protowire.Number, s string) Message {
   for _, str := range []string{s, ""} {
      val, ok := m[Tag{num, str}].(Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number, s string) []byte {
   for _, str := range []string{s, ""} {
      val, ok := m[Tag{num, str}].([]byte)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetMessages(num protowire.Number, s string) []Message {
   for _, str := range []string{s, ""} {
      val, ok := m[Tag{num, str}].([]Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetString(num protowire.Number, s string) string {
   for _, str := range []string{s, ""} {
      val, ok := m[Tag{num, str}].(string)
      if ok {
         return val
      }
   }
   return ""
}

func (m Message) GetUint64(num protowire.Number, s string) uint64 {
   for _, str := range []string{s, ""} {
      val, ok := m[Tag{num, str}].(uint64)
      if ok {
         return val
      }
   }
   return 0
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
         m.addBytes(num, val)
      } else {
         m.addString(num, string(val))
      }
   } else if ok {
      // Could be Message or []byte
      m.Add(num, mes)
   } else {
      // Cound be Message or string
      m.addString(num, string(val))
   }
   return nil
}

func (m Message) consumeFixed32(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint32(num, val)
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, val)
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
   m.Add(num, mes)
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, val)
   return nil
}

////////////////////////////////////////////////////////////////////////////////

func (m Message) Add(num protowire.Number, mes Message) {
   tag := Tag{Number: num}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = mes
   case Message:
      m[tag] = []Message{val, mes}
   case []Message:
      m[tag] = append(val, mes)
   }
}

func (m Message) addBytes(num protowire.Number, v []byte) {
   tag := Tag{Number: num}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case []byte:
      m[tag] = [][]byte{val, v}
   case [][]byte:
      m[tag] = append(val, v)
   }
}

func (m Message) addString(num protowire.Number, v string) {
   tag := Tag{Number: num}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case string:
      m[tag] = []string{val, v}
   case []string:
      m[tag] = append(val, v)
   }
}

func (m Message) addUint32(num protowire.Number, v uint32) {
   tag := Tag{Number: num}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint32:
      m[tag] = []uint32{val, v}
   case []uint32:
      m[tag] = append(val, v)
   }
}

func (m Message) addUint64(num protowire.Number, v uint64) {
   tag := Tag{Number: num}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint64:
      m[tag] = []uint64{val, v}
   case []uint64:
      m[tag] = append(val, v)
   }
}
