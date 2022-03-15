package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

const (
   bytesType = protowire.BytesType
   fixed64Type = protowire.Fixed64Type
   messageType = 6
   varintType = protowire.VarintType
)

type token struct {
   protowire.Number
   protowire.Type
   value interface{}
}

type message []token

func (m message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   for i, tok := range m {
      if tok.Number == num {
         switch value := tok.value.(type) {
         case uint64:
            m[i].value = []uint64{value, val}
         case []uint64:
            m[i].value = append(value, val)
         }
         return nil
      }
   }
   m = append(m, token{num, fixed64Type, val})
   return nil
}

func (m message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   for i, tok := range m {
      if tok.Number == num {
         switch value := tok.value.(type) {
         case uint64:
            m[i].value = []uint64{value, val}
         case []uint64:
            m[i].value = append(value, val)
         }
         return nil
      }
   }
   m = append(m, token{num, varintType, val})
   return nil
}

func (m message) add(num protowire.Number, val message) {
   for i, tok := range m {
      if tok.Number == num && tok.Type == messageType {
         switch value := tok.value.(type) {
         case message:
            m[i].value = []message{value, val}
         case []message:
            m[i].value = append(value, val)
         }
         return
      }
   }
   m = append(m, token{num, messageType, val})
}

func (m message) addString(num protowire.Number, val string) {
   for i, tok := range m {
      if tok.Number == num && tok.Type == bytesType {
         switch value := tok.value.(type) {
         case string:
            m[i].value = []string{value, val}
         case []string:
            m[i].value = append(value, val)
         }
         return
      }
   }
   m = append(m, token{num, bytesType, val})
}

func (m message) addBytes(num protowire.Number, val []byte) {
   for i, tok := range m {
      if tok.Number == num && tok.Type == bytesType {
         switch value := tok.value.(type) {
         case []byte:
            m[i].value = [][]byte{value, val}
         case [][]byte:
            m[i].value = append(value, val)
         }
         return
      }
   }
   m = append(m, token{num, bytesType, val})
}

func (m message) consumeBytes(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := unmarshal(val)
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

func unmarshal(buf []byte) (message, error) {
   var mes message
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
