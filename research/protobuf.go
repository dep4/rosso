package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func add[T any](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = []T{value, val}
   case []T:
      mes[num] = append(value, val)
   }
}

func appendField(buf []byte, num Number, val any) []byte {
   switch val := val.(type) {
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case uint32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case []Message:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []string:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []uint32:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []uint64:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}

func value[T any](mes Message, num Number) T {
   value, _ := mes[num].(T)
   return value
}

func values[T any](mes Message, num Number) []T {
   switch value := mes[num].(type) {
   case []T:
      return value
   case T:
      return []T{value}
   }
   return nil
}

type Message map[Number]any

func Decode(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

func Unmarshal(buf []byte) (Message, error) {
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      if err := protowire.ParseError(fLen); err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      val := buf[tLen:fLen]
      switch typ {
      case protowire.VarintType:
         value, vLen := protowire.ConsumeVarint(val)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add[uint64](mes, num, value)
      case protowire.Fixed64Type:
         value, vLen := protowire.ConsumeFixed64(val)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add[uint64](mes, num, value)
      case protowire.Fixed32Type:
         if err := mes.consumeFixed32(num, val); err != nil {
            return nil, err
         }
      case protowire.BytesType:
         if err := mes.consumeBytes(num, val); err != nil {
            return nil, err
         }
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Add(num Number, val Message) {
   add[Message](m, num, val)
}

func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Message:
      return value
   case string:
      return m.Get(-num)
   }
   return nil
}

func (m Message) GetMessages(num Number) []Message {
   return values[Message](m, num)
}

func (m Message) GetString(num Number) string {
   return value[string](m, num)
}

func (m Message) GetUint64(num Number) uint64 {
   return value[uint64](m, num)
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, value := range m {
      if num >= protowire.MinValidNumber {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}

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

type Number = protowire.Number
