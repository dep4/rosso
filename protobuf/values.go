package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func consumeField(buf []byte) (protowire.Number, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(buf)
   return num, typ, fLen, protowire.ParseError(fLen)
}

func consumeTag(buf []byte) (int, error) {
   _, _, tLen := protowire.ConsumeTag(buf)
   return tLen, protowire.ParseError(tLen)
}

func (m Message) Add(num protowire.Number, str string, val Message) {
   tag := Tag{num, str}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{typ, val}
   case []Message:
      m[tag] = append(typ, val)
   }
}

func (m Message) Get(num protowire.Number, str string) Message {
   val, ok := m[Tag{num, str}].(Message)
   if ok {
      return val
   } else {
      val, ok := m[Tag{num, ""}].(Message)
      if ok {
         return val
      }
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number, str string) []byte {
   val, ok := m[Tag{num, str}].([]byte)
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

func (m Message) GetMessages(num protowire.Number, str string) []Message {
   val, ok := m[Tag{num, str}].([]Message)
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

func (m Message) GetString(num protowire.Number, str string) string {
   val, ok := m[Tag{num, str}].(string)
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

func (m Message) GetUint64(num protowire.Number, str string) uint64 {
   val, ok := m[Tag{num, str}].(uint64)
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

func (m Message) addBytes(num protowire.Number, str string, val []byte) {
   tag := Tag{num, str}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case []byte:
      m[tag] = [][]byte{typ, val}
   case [][]byte:
      m[tag] = append(typ, val)
   }
}

func (m Message) addString(num protowire.Number, str, val string) {
   tag := Tag{num, str}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{typ, val}
   case []string:
      m[tag] = append(typ, val)
   }
}

func (m Message) addUint32(num protowire.Number, str string, val uint32) {
   tag := Tag{num, str}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint32:
      m[tag] = []uint32{typ, val}
   case []uint32:
      m[tag] = append(typ, val)
   }
}

func (m Message) addUint64(num protowire.Number, str string, val uint64) {
   tag := Tag{num, str}
   switch typ := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{typ, val}
   case []uint64:
      m[tag] = append(typ, val)
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
         m.addBytes(num, "", val)
      } else {
         m.addString(num, "", string(val))
      }
   } else if ok {
      // Could be Message or []byte
      m.Add(num, "", mes)
   } else {
      // Cound be Message or string
      m.addString(num, "", string(val))
   }
   return nil
}

func (m Message) consumeFixed32(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint32(num, "", val)
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, "", val)
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
   m.Add(num, "", mes)
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   m.addUint64(num, "", val)
   return nil
}
