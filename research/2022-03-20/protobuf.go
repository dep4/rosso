package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

type Fixed64 uint64

type Fixed32 uint32

type Message map[protowire.Number]any

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
      val := buf[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         err = mes.consumeBytes(num, val)
      case protowire.Fixed64Type:
         err = mes.consumeFixed64(num, val)
      case protowire.Fixed32Type:
         err = mes.consumeFixed32(num, val)
      case protowire.VarintType:
         err = mes.consumeVarint(num, val)
      }
      if err != nil {
         return nil, err
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Add(num protowire.Number, val Message) {
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case Message:
      m[num] = []Message{value, val}
   case []Message:
      m[num] = append(value, val)
   }
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
   if len(val) == 0 {
      return nil
   }
   mes, err := Unmarshal(val)
   if err != nil {
      m.addString(num, string(val))
   } else {
      m.Add(num, mes)
      if !format.IsBinary(val) {
         m.addString(-num, string(val))
      }
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

type Varint uint64
