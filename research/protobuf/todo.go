package protobuf

import (
   "bufio"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func readMessage(r *bufio.Reader) (Message, error) {
   num, typ, tLen := protowire.ConsumeTag(data)
   err := protowire.ParseError(tLen)
   if err != nil {
      return err
   }
   data = data[tLen:]
   var vLen int
   switch typ {
   case protowire.VarintType:
      var val uint64
      val, vLen = protowire.ConsumeVarint(data)
      add(m, num, Varint(val))
   case protowire.Fixed64Type:
      var val uint64
      val, vLen = protowire.ConsumeFixed64(data)
      add(m, num, Fixed64(val))
   case protowire.Fixed32Type:
      var val uint32
      val, vLen = protowire.ConsumeFixed32(data)
      add(m, num, Fixed32(val))
   case protowire.BytesType:
      var val Bytes
      val.Message = make(Message)
      val.Raw, vLen = protowire.ConsumeBytes(data)
      err := val.Message.UnmarshalBinary(val.Raw)
      if err != nil {
         val.Message = nil
      }
      add(m, num, val)
   case protowire.StartGroupType:
      var val Bytes
      val.Message = make(Message)
      val.Raw, vLen = protowire.ConsumeGroup(num, data)
      err := val.Message.UnmarshalBinary(val.Raw)
      if err != nil {
         return err
      }
      add(m, num, val.Message)
   }
}
