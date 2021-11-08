package protobuf

import (
   "encoding/json"
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
   case protowire.BytesType:
      val, vLen := protowire.ConsumeBytes(buf)
      sub := newDecoder(val)
      if sub != nil {
         return sub, vLen
      }
      return string(val), vLen
   case protowire.StartGroupType:
      val, vLen := protowire.ConsumeGroup(num, buf)
      sub := newDecoder(val)
      if sub != nil {
         return sub, vLen
      }
      return val, vLen
   }
   return nil, 0
}

func decode(dec decoder, val interface{}) error {
   buf, err := json.Marshal(dec)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, val)
}

type decoder = map[protowire.Number]interface{}

func newDecoder(buf []byte) decoder {
   dec := make(decoder)
   for len(buf) > 0 {
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
      dVal, ok := dec[num]
      if ok {
         sVal, ok := dVal.([]interface{})
         if ok {
            dec[num] = append(sVal, val)
         } else {
            dec[num] = []interface{}{dVal, val}
         }
      } else {
         dec[num] = val
      }
      buf = buf[fLen:]
   }
   return dec
}
