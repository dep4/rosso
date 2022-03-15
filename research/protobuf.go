package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

// We cannot include the name in the key. When you Unmarshal, the name will be
// empty. If you then try to Get with a name, it will fail. Max valid number is
// 536,870,911, so better to use float64:
// stackoverflow.com/questions/3793838
type Message map[Number]interface{}

func Unmarshal(buf []byte) (Message, error) {
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen, err := consumeField(buf)
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

func (m Message) Add(num Number, name string, val Message) error {
   num += messageType
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case Message:
      m[num] = []Message{value, val}
   case []Message:
      m[num] = append(value, val)
   }
   return nil
}

type Number float64

const (
   messageType Number = 0
   bytesType Number = 0.1
   varintType Number = 0.2
   fixed64Type Number = 0.3
)

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
