package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func consume(num protowire.Number, typ protowire.Type, buf []byte) (interface{}, int) {
   switch typ {
   case protowire.Fixed32Type:
      return protowire.ConsumeFixed32(buf)
   case protowire.Fixed64Type:
      return protowire.ConsumeFixed64(buf)
   case protowire.VarintType:
      return protowire.ConsumeVarint(buf)
   case protowire.StartGroupType:
      buf, vLen := protowire.ConsumeGroup(num, buf)
      recs := unmarshal(buf)
      if recs != nil {
         return recs, vLen
      }
      return buf, vLen
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      if !isBinary(buf) {
         return string(buf), vLen
      }
      recs := unmarshal(buf)
      if recs != nil {
         return recs, vLen
      }
      return buf, vLen
   }
   return nil, 0
}

func isBinary(buf []byte) bool {
   for _, b := range buf {
      switch {
      case b <= 0x08,
      b == 0x0B,
      0x0E <= b && b <= 0x1A,
      0x1C <= b && b <= 0x1F:
         return true
      }
   }
   return false
}

type message map[protowire.Number]interface{}

func unmarshal(buf []byte) message {
   mes := make(message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if tLen <= 0 {
         return nil
      }
      val, vLen := consume(num, typ, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      vMes, ok := mes[num]
      if ok {
         vSlice, ok := vMes.([]interface{})
         if ok {
            mes[num] = append(
               vSlice, token{typ, val},
            )
         } else {
            mes[num] = []interface{}{
               vMes, token{typ, val},
            }
         }
      } else {
         mes[num] = token{typ, val}
      }
      buf = buf[fLen:]
   }
   return mes
}

type token struct {
   Type protowire.Type
   Value interface{}
}
