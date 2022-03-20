package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(buf []byte, num protowire.Number, val any) []byte {
   switch val := val.(type) {
   case Fixed32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, uint32(val))
   case []Fixed32:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case Fixed64:
      buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
      buf = protowire.AppendFixed64(buf, uint64(val))
   case []Fixed64:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case []Message:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []string:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case Varint:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, uint64(val))
   case []Varint:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}

func (m Message) addString(num protowire.Number, val string) {
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case string:
      m[num] = []string{value, val}
   case []string:
      m[num] = append(value, val)
   }
}

func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   if len(val) >= 1 {
      mes, err := Unmarshal(val)
      if err != nil {
         m.addString(num, string(val))
      } else {
         m.Add(num, mes)
         if !format.IsBinary(val) {
            m.addString(-num, string(val))
         }
      }
   } else {
      m.addString(num, "")
   }
   return nil
}

func (m Message) consumeFixed32(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   switch value := m[num].(type) {
   case nil:
      m[num] = Fixed32(val)
   case Fixed32:
      m[num] = []Fixed32{value, Fixed32(val)}
   case []Fixed32:
      m[num] = append(value, Fixed32(val))
   }
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   switch value := m[num].(type) {
   case nil:
      m[num] = Fixed64(val)
   case Fixed64:
      m[num] = []Fixed64{value, Fixed64(val)}
   case []Fixed64:
      m[num] = append(value, Fixed64(val))
   }
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   switch value := m[num].(type) {
   case nil:
      m[num] = Varint(val)
   case Varint:
      m[num] = []Varint{value, Varint(val)}
   case []Varint:
      m[num] = append(value, Varint(val))
   }
   return nil
}
