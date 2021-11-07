package protobuf

import (
   "bytes"
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

func appendField(out []byte, key protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case Message:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendBytes(out, val.Marshal())
   case Repeated:
      for _, v := range val {
         out = appendField(out, key, v)
      }
   case bool:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, protowire.EncodeBool(val))
   case int32:
      out = protowire.AppendTag(out, key, protowire.VarintType)
      out = protowire.AppendVarint(out, uint64(val))
   case string:
      out = protowire.AppendTag(out, key, protowire.BytesType)
      out = protowire.AppendString(out, val)
   }
   return out
}

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

type Message map[protowire.Number]interface{}

func Parse(buf []byte) Message {
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

func (m Message) Marshal() []byte {
   var out []byte
   for key, val := range m {
      out = appendField(out, key, val)
   }
   return out
}

func (m Message) Transcode(v interface{}) error {
   buf := new(bytes.Buffer)
   err := json.NewEncoder(buf).Encode(m)
   if err != nil {
      return err
   }
   return json.NewDecoder(buf).Decode(v)
}

type Repeated []interface{}
