package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func add[T Encoder](m map[Number]Encoder, num Number, val T) bool {
   vals, ok := m[num].(T)
   if ok {
      m[num] = append(vals, val)
   }
   return ok
}

func Unmarshal(buf []byte) (map[Number]Encoder, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(map[Number]Encoder)
   for len(buf) >= 1 {
      //num, typ, tLen := protowire.ConsumeTag(buf)
      _, typ, tLen := protowire.ConsumeTag(buf)
      err := protowire.ParseError(tLen)
      if err != nil {
         return nil, err
      }
      buf = buf[tLen:]
      var vLen int
      switch typ {
      case protowire.VarintType:
         //var val uint64
         //val, vLen = protowire.ConsumeVarint(buf)
         //mes[num] = append(mes[num], Varint(val))
      case protowire.Fixed32Type:
         //var val uint32
         //val, vLen = protowire.ConsumeFixed32(buf)
         //mes[num] = append(mes[num], Fixed32(val))
      case protowire.Fixed64Type:
         //var val uint64
         //val, vLen = protowire.ConsumeFixed64(buf)
         //mes[num] = append(mes[num], Fixed64(val))
      case protowire.BytesType:
         /*
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeBytes(buf)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         mes[num] = append(mes[num], val)
         */
      }
      if err := protowire.ParseError(vLen); err != nil {
         return nil, err
      }
      buf = buf[vLen:]
   }
   return mes, nil
}
