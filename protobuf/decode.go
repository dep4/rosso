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
      dec := NewDecoder(buf)
      if dec != nil {
         return dec, vLen
      }
      return buf, vLen
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      if isText(buf) {
         return string(buf), vLen
      }
      dec := NewDecoder(buf)
      if dec != nil {
         return dec, vLen
      }
      return buf, vLen
   }
   return nil, 0
}

// github.com/golang/go/blob/go1.17.3/src/net/http/sniff.go#L297-L309
func isText(buf []byte) bool {
   for _, b := range buf {
      switch {
      case b <= 0x08,
      b == 0x0B,
      0x0E <= b && b <= 0x1A,
      0x1C <= b && b <= 0x1F:
         return false
      }
   }
   return true
}

type Decoder map[protowire.Number]interface{}

// Convert byte slice to map
func NewDecoder(buf []byte) Decoder {
   dec := make(Decoder)
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

// Convert map to struct
func (d Decoder) Decode(val interface{}) error {
   buf, err := json.Marshal(d)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, val)
}
