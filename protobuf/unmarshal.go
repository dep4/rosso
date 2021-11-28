package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func Unmarshal(buf []byte) Message {
   mes := make(Message)
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
         /*
         switch vMes.(type) {
         case uint32:
         case uint64:
         case Bytes:
         case Message:
         case []uint32:
         case []uint64:
         case []Bytes:
         case []Message:
         }
         */
         vSlice, ok := vMes.([]interface{})
         if ok {
            mes[num] = append(vSlice, val)
         } else {
            mes[num] = []interface{}{vMes, val}
         }
      } else {
         mes[num] = val
      }
      buf = buf[fLen:]
   }
   return mes
}
