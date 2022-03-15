package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

const (
   bytesType = 2
   fixed64Type = 1
   messageType = 6
   varintType = 0
)

type Message []Field

func Unmarshal(buf []byte) (Message, error) {
   var mes Message
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
      case protowire.Fixed64Type:
         err = mes.consumeFixed64(num, val)
      case protowire.VarintType:
         err = mes.consumeVarint(num, val)
      case protowire.BytesType:
         err = mes.consumeBytes(num, val)
      }
      if err != nil {
         return nil, err
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) add(num protowire.Number, val Message) {
   for i, field := range m {
      if field.Number == num && field.Type == messageType {
         switch value := field.Value.(type) {
         case Message:
            m[i].Value = []Message{value, val}
         case []Message:
            m[i].Value = append(value, val)
         }
         return
      }
   }
   m = append(m, Field{num, messageType, val})
}

func (m Message) addBytes(num protowire.Number, val []byte) {
   for i, field := range m {
      if field.Number == num && field.Type == bytesType {
         switch value := field.Value.(type) {
         case []byte:
            m[i].Value = [][]byte{value, val}
         case [][]byte:
            m[i].Value = append(value, val)
         }
         return
      }
   }
   m = append(m, Field{num, bytesType, val})
}

func (m Message) addString(num protowire.Number, val string) {
   for i, field := range m {
      if field.Number == num && field.Type == bytesType {
         switch value := field.Value.(type) {
         case string:
            m[i].Value = []string{value, val}
         case []string:
            m[i].Value = append(value, val)
         }
         return
      }
   }
   m = append(m, Field{num, bytesType, val})
}

func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         m.addBytes(num, val)
      } else {
         m.addString(num, string(val))
      }
   } else {
      m.add(num, mes)
      if !binary {
         m.addString(num, string(val))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   for i, field := range m {
      if field.Number == num {
         switch value := field.Value.(type) {
         case uint64:
            m[i].Value = []uint64{value, val}
         case []uint64:
            m[i].Value = append(value, val)
         }
         return nil
      }
   }
   m = append(m, Field{num, fixed64Type, val})
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   for i, field := range m {
      if field.Number == num {
         switch value := field.Value.(type) {
         case uint64:
            m[i].Value = []uint64{value, val}
         case []uint64:
            m[i].Value = append(value, val)
         }
         return nil
      }
   }
   m = append(m, Field{num, varintType, val})
   return nil
}

type Field struct {
   protowire.Number
   protowire.Type
   Value interface{}
}
