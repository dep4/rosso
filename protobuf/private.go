package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

const (
   // use messageType as the default Type for Tag
   messageType = 0
   fixed64Type = 1
   bytesType = 2
   varintType = 6
)

func appendField(b []byte, num protowire.Number, val any) []byte {
   switch val := val.(type) {
   case uint64:
      b = protowire.AppendTag(b, num, protowire.VarintType)
      b = protowire.AppendVarint(b, val)
   case string:
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendString(b, val)
   case Message:
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendBytes(b, val.Marshal())
   case []uint64:
      for _, value := range val {
         b = appendField(b, num, value)
      }
   case []string:
      for _, value := range val {
         b = appendField(b, num, value)
      }
   case []Message:
      for _, value := range val {
         b = appendField(b, num, value)
      }
   }
   return b
}

func (m Message) addString(num protowire.Number, val string) {
   key := Tag{num, bytesType}
   switch value := m[key].(type) {
   case nil:
      m[key] = val
   case string:
      m[key] = []string{value, val}
   case []string:
      m[key] = append(value, val)
   }
}

func (m Message) consumeBytes(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeBytes(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   mes, err := Unmarshal(val)
   if err != nil {
      m.addString(num, string(val))
   } else {
      m.Add(num, "", mes)
      if !format.IsBinary(val) {
         m.addString(num, string(val))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   key := Tag{num, fixed64Type}
   switch value := m[key].(type) {
   case nil:
      m[key] = val
   case uint64:
      m[key] = []uint64{value, val}
   case []uint64:
      m[key] = append(value, val)
   }
   return nil
}

func (m Message) consumeVarint(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeVarint(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   key := Tag{num, varintType}
   switch value := m[key].(type) {
   case nil:
      m[key] = val
   case uint64:
      m[key] = []uint64{value, val}
   case []uint64:
      m[key] = append(value, val)
   }
   return nil
}
