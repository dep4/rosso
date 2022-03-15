package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m message) consumeVarint(num protowire.Number, buf []byte) error {
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

func (m message) Add(num Number, name string, val message) error {
   num += messageType
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case message:
      m[num] = []message{value, val}
   case []message:
      m[num] = append(value, val)
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

func (m message) consumeBytes(num Number, buf []byte) error {
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

func (m message) addString(num Number, val string) {
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
