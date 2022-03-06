package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

const (
   messageType Number = 0
   bytesType Number = 0.1
   varintType Number = 0.2
   fixed64Type Number = 0.3
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

// In some cases if input is binary, then result could be a Message or byte
// slice. We assume for now its always a Message. If input is not binary, then
// result could be a Message or string. Since its not possible to tell Message
// from string, we just add both under the same number, each with its own type.
func (m Message) consumeBytes(num Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         num += bytesType
         switch value := m[num].(type) {
         case nil:
            m[num] = val
         case []byte:
            m[num] = [][]byte{value, val}
         case [][]byte:
            m[num] = append(value, val)
         }
      } else {
         m.addString(num, string(val))
      }
   } else {
      m.Add(num, "", mes)
      if !binary {
         m.addString(num, string(val))
      }
   }
   return nil
}

func consumeField(buf []byte) (Number, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(buf)
   return Number(num), typ, fLen, protowire.ParseError(fLen)
}

func (m Message) addString(num Number, val string) {
   num += bytesType
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case string:
      m[num] = []string{value, val}
   case []string:
      m[num] = append(value, val)
   }
}

func (m Message) consumeFixed64(num Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   num += fixed64Type
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
   num += varintType
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
