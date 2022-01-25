package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) Add(num protowire.Number, s string, v Message) {
   tag := Tag{num, s}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = v
   case Message:
      m[tag] = []Message{typ, v}
   case []Message:
      m[tag] = append(typ, v)
   }
}

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
   val, ok := m[Tag{num, s}].([]byte)
   if ok {
      return val
   } else {
      val, ok := m[Tag{num, ""}].([]byte)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetMessages(num protowire.Number, s string) []Message {
   val, ok := m[Tag{num, s}].([]Message)
   if ok {
      return val
   } else {
      val, ok := m[Tag{num, ""}].([]Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetString(num protowire.Number, s string) string {
   val, ok := m[Tag{num, s}].(string)
   if ok {
      return val
   } else {
      val, ok := m[Tag{num, ""}].(string)
      if ok {
         return val
      }
   }
   return ""
}

func (m Message) GetUint64(num protowire.Number, s string) uint64 {
   val, ok := m[Tag{num, s}].(uint64)
   if ok {
      return val
   } else {
      val, ok := m[Tag{num, ""}].(uint64)
      if ok {
         return val
      }
   }
   return 0
}

func (m Message) addBytes(num protowire.Number, s string, v []byte) {
   tag := Tag{num, s}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = v
   case []byte:
      m[tag] = [][]byte{typ, v}
   case [][]byte:
      m[tag] = append(typ, v)
   }
}

func (m Message) addString(num protowire.Number, s, v string) {
   tag := Tag{num, s}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = v
   case string:
      m[tag] = []string{typ, v}
   case []string:
      m[tag] = append(typ, v)
   }
}

func (m Message) addUint32(num protowire.Number, s string, v uint32) {
   tag := Tag{num, s}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint32:
      m[tag] = []uint32{typ, v}
   case []uint32:
      m[tag] = append(typ, v)
   }
}

func (m Message) addUint64(num protowire.Number, s string, v uint64) {
   tag := Tag{num, s}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = v
   case uint64:
      m[tag] = []uint64{typ, v}
   case []uint64:
      m[tag] = append(typ, v)
   }
}

////////////////////////////////////////////////////////////////////////////////

func (m Message) consumeBytes(num protowire.Number, b []byte) error {
   buf, vLen := protowire.ConsumeBytes(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   ok := format.IsBinary(buf)
   mes, err := Unmarshal(buf)
   if err != nil {
      if ok {
         m.addBytes(num, "", buf)
      } else {
         m.addString(num, "", string(buf))
      }
   } else if ok {
      // Could be Message or []byte
      m.Add(num, "", mes)
   } else {
      // Cound be Message or string
      m.addString(num, "", string(buf))
   }
   return nil
}

func (m Message) consumeFixed32(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint32(num, "", val)
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, "", val)
   return nil
}

func (m Message) consumeGroup(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeGroup(num, b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   mes, err := Unmarshal(val)
   if err != nil {
      return err
   }
   m.Add(num, "", mes)
   return nil
}

func (m Message) consumeVarint(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeVarint(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, "", val)
   return nil
}
