package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

const (
   messageType = 0
   fixed64Type = 1
   bytesType = 2
   varintType = 6
)

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

type Message map[Tag]any

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

func (m Message) Add(num protowire.Number, s string, v Message) {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = v
   case Message:
      m[tag] = []Message{value, v}
   case []Message:
      m[tag] = append(value, v)
   }
}

type Tag struct {
   protowire.Number
   protowire.Type
}
