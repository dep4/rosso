package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func consumeBytes(b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeBytes(b)
   return val, protowire.ParseError(vLen)
}

func consumeField(b []byte) (protowire.Number, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(b)
   return num, typ, fLen, protowire.ParseError(fLen)
}

func consumeFixed32(b []byte) (uint32, error) {
   val, vLen := protowire.ConsumeFixed32(b)
   return val, protowire.ParseError(vLen)
}

func consumeFixed64(b []byte) (uint64, error) {
   val, vLen := protowire.ConsumeFixed64(b)
   return val, protowire.ParseError(vLen)
}

func consumeGroup(num protowire.Number, b []byte) ([]byte, error) {
   val, vLen := protowire.ConsumeGroup(num, b)
   return val, protowire.ParseError(vLen)
}

func consumeTag(b []byte) (int, error) {
   _, _, tLen := protowire.ConsumeTag(b)
   return tLen, protowire.ParseError(tLen)
}

func consumeVarint(b []byte) (uint64, error) {
   val, vLen := protowire.ConsumeVarint(b)
   return val, protowire.ParseError(vLen)
}

func (m Message) bytesType(num protowire.Number, buf []byte) {
   ok := format.IsBinary(buf)
   mNew, err := Unmarshal(buf)
   if err != nil {
      if ok {
         m.addBytes(num, "", buf)
      } else {
         m.addString(num, "", string(buf))
      }
   } else if ok {
      // Could be Message or []byte
      m.Add(num, "", mNew)
   } else {
      // Cound be Message or string
      m.addString(num, "", string(buf))
   }
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

func (m Message) fixed32Type(num protowire.Number, val uint32) {
   tag := Tag{num, ""}
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
