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
   case protowire.StartGroupType:
      buf, vLen := protowire.ConsumeGroup(num, buf)
      recs := Bytes(buf)
      if recs != nil {
         return recs, vLen
      }
      return buf, vLen
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      if ! isBinary(buf) {
         return string(buf), vLen
      }
      recs := Bytes(buf)
      if recs != nil {
         return recs, vLen
      }
      return buf, vLen
   }
   return nil, 0
}

// mimesniff.spec.whatwg.org#binary-data-byte
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

type Records map[protowire.Number]interface{}

func Bytes(buf []byte) Records {
   recs := make(Records)
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
      dVal, ok := recs[num]
      if ok {
         sVal, ok := dVal.([]interface{})
         if ok {
            recs[num] = append(sVal, val)
         } else {
            recs[num] = []interface{}{dVal, val}
         }
      } else {
         recs[num] = val
      }
      buf = buf[fLen:]
   }
   return recs
}

func (r Records) Struct(val interface{}) error {
   buf, err := json.Marshal(r)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, val)
}
