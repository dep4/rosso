package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

func consume(key protowire.Number, typ protowire.Type, buf []byte) (interface{}, int) {
   switch typ {
   case protowire.Fixed32Type:
      return protowire.ConsumeFixed32(buf)
   case protowire.Fixed64Type:
      return protowire.ConsumeFixed64(buf)
   case protowire.VarintType:
      return protowire.ConsumeVarint(buf)
   case protowire.BytesType:
      v, vLen := protowire.ConsumeBytes(buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return string(v), vLen
   case protowire.StartGroupType:
      v, vLen := protowire.ConsumeGroup(key, buf)
      sub := Parse(v)
      if sub != nil {
         return sub, vLen
      }
      return v, vLen
   }
   return nil, 0
}

func alfa(buf []byte) map[protowire.Number]interface{} {
   mes := make(Message)
   for len(buf) > 0 {
      key, typ, fLen := protowire.ConsumeField(buf)
      if fLen <= 0 {
         return nil
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if tLen <= 0 {
         return nil
      }
      v, vLen := consume(key, typ, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
      iface, ok := mes[key]
      if ok {
         rep, ok := iface.(Repeated)
         if ok {
            mes[key] = append(rep, v)
         } else {
            mes[key] = Repeated{iface, v}
         }
      } else {
         mes[key] = v
      }
      buf = buf[fLen:]
   }
   return mes
}

func bravo(m map[]interface{}, v interface{}) error {
   buf, err := json.Marshal(m)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, any)
}
