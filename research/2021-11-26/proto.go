package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case token:
      buf = protowire.AppendTag(buf, num, val.Type)
      switch val.Type {
      case protowire.Fixed32Type:
         buf = protowire.AppendFixed32(buf, val.Value.(uint32))
      case protowire.Fixed64Type:
         buf = protowire.AppendFixed64(buf, val.Value.(uint64))
      case protowire.VarintType:
         buf = protowire.AppendVarint(buf, val.Value.(uint64))
      case protowire.BytesType, protowire.StartGroupType:
         switch val := val.Value.(type) {
         case string:
            buf = protowire.AppendString(buf, val)
         case []byte:
            buf = protowire.AppendBytes(buf, val)
         case message:
            buf = protowire.AppendBytes(buf, val.marshal())
         }
      }
   case []interface{}:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   }
   return buf
}

func consume(num protowire.Number, typ protowire.Type, buf []byte) (interface{}, int) {
   switch typ {
   case protowire.VarintType:
      return protowire.ConsumeVarint(buf)
   case protowire.Fixed32Type:
      return protowire.ConsumeFixed32(buf)
   case protowire.Fixed64Type:
      return protowire.ConsumeFixed64(buf)
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      recs := unmarshal(buf)
      if recs != nil {
         return recs, vLen
      }
      if isBinary(buf) {
         return buf, vLen
      }
      return string(buf), vLen
   case protowire.StartGroupType:
      buf, vLen := protowire.ConsumeGroup(num, buf)
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
      tok := token{typ, val}
      vMes, ok := mes[num]
      if ok {
         vSlice, ok := vMes.([]interface{})
         if ok {
            mes[num] = append(vSlice, tok)
         } else {
            mes[num] = []interface{}{vMes, tok}
         }
      } else {
         mes[num] = tok
      }
      buf = buf[fLen:]
   }
   return mes
}

func (m message) marshal() []byte {
   var buf []byte
   for key, val := range m {
      buf = appendField(buf, key, val)
   }
   return buf
}

type token struct {
   Type protowire.Type
   Value interface{}
}
