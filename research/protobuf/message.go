package protobuf

import (
   "errors"
   "github.com/89z/format"
   "io"
   "google.golang.org/protobuf/encoding/protowire"
)

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, length := protowire.ConsumeTag(buf)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      buf = buf[length:]
      switch typ {
      case protowire.VarintType:
         buf, err = mes.consume_varint(num, buf)
      case protowire.Fixed64Type:
         buf, err = mes.consume_fixed64(num, buf)
      case protowire.Fixed32Type:
         buf, err = mes.consume_fixed32(num, buf)
      case protowire.BytesType:
         buf, err = mes.consume_raw(num, buf)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}

func (m Message) consume_raw(num Number, b []byte) ([]byte, error) {
   var (
      length int
      val Raw
   )
   val.Bytes, length = protowire.ConsumeBytes(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if format.String(val.Bytes) {
      val.String = string(val.Bytes)
   }
   val.Message, _ = Unmarshal(val.Bytes)
   if err := add(m, num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

type Message map[Number]Encoder

func (m Message) Fixed64(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Fixed64)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) Varint(num Number) (uint64, error) {
   lvalue := m[num]
   rvalue, ok := lvalue.(Varint)
   if !ok {
      return 0, type_error{num, lvalue, rvalue}
   }
   return uint64(rvalue), nil
}

func (m Message) consume_fixed32(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed32(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Fixed32(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (m Message) Add_Fixed32(num Number, v uint32) error {
   return add(m, num, Fixed32(v))
}

func (m Message) Add_Fixed64(num Number, v uint64) error {
   return add(m, num, Fixed64(v))
}

func (m Message) consume_fixed64(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed64(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Fixed64(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

func (m Message) Add_Varint(num Number, v uint64) error {
   return add(m, num, Varint(v))
}

func (m Message) consume_varint(num Number, b []byte) ([]byte, error) {
   val, length := protowire.ConsumeVarint(b)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := m.Add_Varint(num, val); err != nil {
      return nil, err
   }
   return b[length:], nil
}

