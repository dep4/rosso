package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
)

func unmarshalJSON(buf []byte) (interface{}, error) {
   return nil, nil
}

////////////////////////////////////////////////////////////////////////////////

func (m *message) UnmarshalJSON(buf []byte) error {
   var raw map[protowire.Number]json.RawMessage
   err := json.Unmarshal(buf, &raw)
   if err != nil {
      return err
   }
   for key, val := range raw {
      any, err := unmarshalJSON(val)
      if err != nil {
         return err
      }
      (*m)[key] = any
   }
   return nil
}

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
      case protowire.BytesType:
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

func consume(num protowire.Number, typ protowire.Type, buf []byte) (token, int) {
   switch typ {
   case protowire.Fixed32Type:
      val, vLen := protowire.ConsumeFixed32(buf)
      return token{typ, val}, vLen
   case protowire.Fixed64Type:
      val, vLen := protowire.ConsumeFixed64(buf)
      return token{typ, val}, vLen
   case protowire.VarintType:
      val, vLen := protowire.ConsumeVarint(buf)
      return token{typ, val}, vLen
   case protowire.StartGroupType:
      buf, vLen := protowire.ConsumeGroup(num, buf)
      recs := unmarshal(buf)
      if recs != nil {
         return token{protowire.BytesType, recs}, vLen
      }
      return token{protowire.BytesType, buf}, vLen
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      if !isBinary(buf) {
         return token{typ, string(buf)}, vLen
      }
      recs := unmarshal(buf)
      if recs != nil {
         return token{protowire.BytesType, recs}, vLen
      }
      return token{protowire.BytesType, buf}, vLen
   }
   return token{}, 0
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
      tok, vLen := consume(num, typ, buf[tLen:fLen])
      if vLen <= 0 {
         return nil
      }
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

func (m message) MarshalJSON() ([]byte, error) {
   mes := map[protowire.Number]interface{}(m)
   return json.Marshal(mes)
}

type token struct {
   Type protowire.Type
   Value interface{}
}
