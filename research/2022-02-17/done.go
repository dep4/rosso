package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Message map[Tag]interface{}

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

const (
   bytesType = iota * 0.1
   fixed64Type
   messageType
   stringType
   varintType
)

type Tag struct {
   NumberType float64
   Name string
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val2, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{NumberType: num + fixed64Type}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = val2
   case uint64:
      m[tag] = []uint64{val, val2}
   case []uint64:
      m[tag] = append(val, val2)
   }
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val2, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{NumberType: num + varintType}
   switch val := m[tag].(type) {
   case nil:
      m[tag] = val2
   case uint64:
      m[tag] = []uint64{val, val2}
   case []uint64:
      m[tag] = append(val, val2)
   }
   return nil
}

// In some cases if input is binary, then result could be a Message or byte
// slice. We assume for now its always a Message. If input is not binary, then
// result could be a Message or string. Since its not possible to tell Message
// from string, we just add both under the same number, each with its own type.
func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   val2, eLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val2)
   mes, err := Unmarshal(val2)
   if err != nil {
      if binary {
         tag := Tag{NumberType: num + bytesType}
         switch val := m[tag].(type) {
         case nil:
            m[tag] = val2
         case []byte:
            m[tag] = [][]byte{val, val2}
         case [][]byte:
            m[tag] = append(val, val2)
         }
      } else {
         m.addString(num, string(val2))
      }
   } else {
      m.Add(num, "", mes)
      if !binary {
         m.addString(num, string(val2))
      }
   }
   return nil
}
