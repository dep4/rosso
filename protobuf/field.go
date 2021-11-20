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
      nmap := NewFields(buf)
      if nmap != nil {
         return nmap, vLen
      }
      return buf, vLen
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      if isText(buf) {
         return string(buf), vLen
      }
      nmap := NewFields(buf)
      if nmap != nil {
         return nmap, vLen
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

type Fields map[protowire.Number]interface{}

// Convert byte slice to map
func NewFields(buf []byte) Fields {
   nmap := make(Fields)
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
      dVal, ok := nmap[num]
      if ok {
         sVal, ok := dVal.([]interface{})
         if ok {
            nmap[num] = append(sVal, val)
         } else {
            nmap[num] = []interface{}{dVal, val}
         }
      } else {
         nmap[num] = val
      }
      buf = buf[fLen:]
   }
   return nmap
}

// Convert map to struct
func (f Fields) Struct(val interface{}) error {
   buf, err := json.Marshal(f)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, val)
}
