package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func (m Message) Add(num Number, val Message) {
   add[Message](m, num, val)
}

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

func values[T any](mes Message, num Number) []T {
   switch value := mes[num].(type) {
   case []T:
      return value
   case T:
      return []T{value}
   }
   return nil
}

func (m Message) GetMessages(num Number) []Message {
   return values[Message](m, num)
}

func value[T any](mes Message, num Number) T {
   value, _ := mes[num].(T)
   return value
}

func (m Message) GetString(num Number) string {
   return value[string](m, num)
}

func (m Message) GetUint64(num Number) uint64 {
   return value[uint64](m, num)
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
      err := protowire.ParseError(fLen)
      if err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      value := buf[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         err = mes.consumeBytes(num, value)
      case protowire.Fixed64Type:
         err = mes.consumeFixed64(num, value)
      case protowire.Fixed32Type:
         err = mes.consumeFixed32(num, value)
      case protowire.VarintType:
         err = mes.consumeVarint(num, value)
      }
      if err != nil {
         return nil, err
      }
      buf = buf[fLen:]
   }
   return mes, nil
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

type Number = protowire.Number
