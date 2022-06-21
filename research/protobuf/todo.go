package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
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
         buf, err = mes.consume_bytes(num, buf)
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}

func (m Message) consume_fixed64(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Fixed64](m, num)
   if err != nil {
      return nil, err
   }
   val, length := protowire.ConsumeFixed64(buf)
   if err := protowire.ParseError(length); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[length:], nil
}

func (m Message) consume_bytes(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Bytes](m, num)
   if err != nil {
      return nil, err
   }
   var (
      val Bytes
      length int
   )
   val.Raw, length = protowire.ConsumeBytes(buf)
   if err := protowire.ParseError(length); err != nil {
      return nil, err
   }
   val.Message, err = Unmarshal(val.Raw)
   if err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[length:], nil
}

func (m Message) consume_fixed32(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Fixed32](m, num)
   if err != nil {
      return nil, err
   }
   val, length := protowire.ConsumeFixed32(buf)
   if err := protowire.ParseError(length); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[length:], nil
}
