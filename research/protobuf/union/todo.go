package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func (m Message) UnmarshalBinary(buf []byte) error {
   if len(buf) == 0 {
      return io.ErrUnexpectedEOF
   }
   for len(buf) >= 1 {
      num, typ, tLen := protowire.ConsumeTag(buf)
      err := protowire.ParseError(tLen)
      if err != nil {
         return err
      }
      buf = buf[tLen:]
      tok := m[num]
      var (
         vLen int
         val Value
      )
      switch typ {
      case protowire.VarintType:
         val.Int64, vLen = protowire.ConsumeVarint(buf)
         tok.Value = append(tok.Value, val)
         m[num] = tok
      case protowire.Fixed64Type:
         val.Int64, vLen = protowire.ConsumeFixed64(buf)
         tok.Value = append(tok.Value, val)
         m[num] = tok
      case protowire.Fixed32Type:
         val.Int32, vLen = protowire.ConsumeFixed32(buf)
         tok.Value = append(tok.Value, val)
         m[num] = tok
      case protowire.BytesType:
         /*
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeBytes(buf)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         add(m, num, val)
         */
      }
      if err := protowire.ParseError(vLen); err != nil {
         return err
      }
      buf = buf[vLen:]
   }
   return nil
}
