package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
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

func (m Message) Add(num Number, val Message) {
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case Message:
      m[num] = []Message{value, val}
   case []Message:
      m[num] = append(value, val)
   }
}

func (m Message) Get(num Number) Message {
   switch val := m[num].(type) {
   case Message:
      return val
   case string:
      return m.Get(-num)
   }
   return nil
}

func (m Message) GetMessages(num Number) []Message {
   switch val := m[num].(type) {
   case []Message:
      return val
   case Message:
      return []Message{val}
   }
   return nil
}

func (m Message) GetString(num Number) string {
   val, _ := m[num].(string)
   return val
}

func (m Message) GetUint64(num Number) uint64 {
   val, _ := m[num].(uint64)
   return val
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, val := range m {
      if num >= protowire.MinValidNumber {
         buf = appendField(buf, num, val)
      }
   }
   return buf
}

type Number = protowire.Number
