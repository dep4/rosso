package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m Message) addString(num Number, val string) {
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case string:
      m[num] = []string{value, val}
   case []string:
      m[num] = append(value, val)
   }
}

func (m Message) consumeBytes(num Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   if len(val) >= 1 {
      mes, err := Unmarshal(val)
      if err != nil {
         m.addString(num, string(val))
      } else if format.IsBinary(val) {
         m.Add(num, mes)
      } else {
         // Message should be negative, as string is easier to Marshal
         m.addString(num, string(val))
         m.Add(-num, mes)
      }
   } else {
      m.addString(num, "")
   }
   return nil
}

func (m Message) consumeFixed32(num Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case uint32:
      m[num] = []uint32{value, val}
   case []uint32:
      m[num] = append(value, val)
   }
   return nil
}

func (m Message) consumeFixed64(num Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case uint64:
      m[num] = []uint64{value, val}
   case []uint64:
      m[num] = append(value, val)
   }
   return nil
}

func (m Message) consumeVarint(num Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case uint64:
      m[num] = []uint64{value, val}
   case []uint64:
      m[num] = append(value, val)
   }
   return nil
}
