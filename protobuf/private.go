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

func appendField(b []byte, num protowire.Number, v interface{}) []byte {
   switch v := v.(type) {
   case uint64:
      b = protowire.AppendTag(b, num, protowire.VarintType)
      b = protowire.AppendVarint(b, v)
   case string:
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendString(b, v)
   case []byte:
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendBytes(b, v)
   case Message:
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendBytes(b, v.Marshal())
   case []uint64:
      for _, value := range v {
         b = appendField(b, num, value)
      }
   case []string:
      for _, value := range v {
         b = appendField(b, num, value)
      }
   case []Message:
      for _, value := range v {
         b = appendField(b, num, value)
      }
   }
   return b
}

func (m Message) addString(num protowire.Number, v string) {
   tag := Tag{num, bytesType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = v
   case string:
      m[tag] = []string{value, v}
   case []string:
      m[tag] = append(value, v)
   }
}

func (m Message) consumeBytes(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeBytes(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         tag := Tag{num, bytesType}
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
      m.Add(num, "", mes)
      if !binary {
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
   tag := Tag{num, fixed64Type}
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

func (m Message) consumeVarint(num protowire.Number, b []byte) error {
   val, vLen := protowire.ConsumeVarint(b)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{num, varintType}
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
