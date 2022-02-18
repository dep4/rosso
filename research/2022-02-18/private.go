package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

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
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []string:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []Message:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}

func consumeField(buf []byte) (float64, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(buf)
   return float64(num), typ, fLen, protowire.ParseError(fLen)
}

func (m Message) addString(num float64, val string) {
   tag := Tag{NumberType: num + BytesType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{value, val}
   case []string:
      m[tag] = append(value, val)
   }
}

// In some cases if input is binary, then result could be a Message or byte
// slice. We assume for now its always a Message. If input is not binary, then
// result could be a Message or string. Since its not possible to tell Message
// from string, we just add both under the same number, each with its own type.
func (m Message) consumeBytes(num float64, buf []byte) error {
   val, eLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         tag := Tag{NumberType: num + BytesType}
         switch value := m[tag].(type) {
         case nil:
            m[tag] = val
         case []byte:
            m[tag] = [][]byte{value, val}
         case [][]byte:
            m[tag] = append(value, val)
         }
      } else {
         m.addString(num, string(val))
      }
   } else {
      m.Add(num, mes)
      if !binary {
         m.addString(num, string(val))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{NumberType: num + Fixed64Type}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{value, val}
   case []uint64:
      m[tag] = append(value, val)
   }
   return nil
}

func (m Message) consumeVarint(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{NumberType: num + VarintType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{value, val}
   case []uint64:
      m[tag] = append(value, val)
   }
   return nil
}
