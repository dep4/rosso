package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

type Message map[Number]any

func Decode(in io.Reader) (Message, error) {
   buf, err := io.ReadAll(in)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

func Unmarshal(in []byte) (Message, error) {
   mes := make(Message)
   for len(in) >= 1 {
      num, typ, fLen := protowire.ConsumeField(in)
      if err := protowire.ParseError(fLen); err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(in[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      buf := in[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         if len(val) >= 1 {
            embed, err := Unmarshal(val)
            if err != nil {
               add(mes, num, string(val))
            } else if format.IsBinary(val) {
               add(mes, num, embed)
            } else {
               add(mes, num, string(val))
               add(mes, -num, embed)
            }
         } else {
            add(mes, num, "")
         }
      case protowire.Fixed32Type:
         val, vLen := protowire.ConsumeFixed32(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(buf)
         if err := protowire.ParseError(vLen); err != nil {
            return nil, err
         }
         add(mes, num, val)
      }
      in = in[fLen:]
   }
   return mes, nil
}

func (m Message) Add(num Number, val Message) {
   add(m, num, val)
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
   switch value := m[num].(type) {
   case []Message:
      return value
   case Message:
      return []Message{value}
   }
   return nil
}

func (m Message) GetString(num Number) string {
   return get[string](m, num)
}

func (m Message) GetUint64(num Number) uint64 {
   return get[uint64](m, num)
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
