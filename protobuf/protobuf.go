package protobuf

import (
   "encoding/json"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case bool:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, protowire.EncodeBool(val))
   case float64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, uint64(val))
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []interface{}:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case map[protowire.Number]interface{}:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, Message.Marshal(val))
   }
   return buf
}

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
      recs := NewMessageBytes(buf)
      if recs != nil {
         return recs, vLen
      }
      return buf, vLen
   case protowire.BytesType:
      buf, vLen := protowire.ConsumeBytes(buf)
      if !isBinary(buf) {
         return string(buf), vLen
      }
      recs := NewMessageBytes(buf)
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

func unmarshal(buf []byte) (interface{}, error) {
   if buf[0] == '{' {
      mes := make(Message)
      err := json.Unmarshal(buf, &mes)
      if err != nil {
         return nil, err
      }
      return mes, nil
   }
   if buf[0] == '[' {
      var raw []json.RawMessage
      err := json.Unmarshal(buf, &raw)
      if err != nil {
         return nil, err
      }
      var arr []interface{}
      for _, val := range raw {
         any, err := unmarshal(val)
         if err != nil {
            return nil, err
         }
         arr = append(arr, any)
      }
      return arr, nil
   }
   var any interface{}
   err := json.Unmarshal(buf, &any)
   if err != nil {
      return nil, err
   }
   return any, nil
}

type Message map[protowire.Number]interface{}

func NewMessage(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return NewMessageBytes(buf), nil
}

func NewMessageBytes(buf []byte) Message {
   mes := make(Message)
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
      dVal, ok := mes[num]
      if ok {
         sVal, ok := dVal.([]interface{})
         if ok {
            mes[num] = append(sVal, val)
         } else {
            mes[num] = []interface{}{dVal, val}
         }
      } else {
         mes[num] = val
      }
      buf = buf[fLen:]
   }
   return mes
}

func (m Message) Marshal() []byte {
   var buf []byte
   for key, val := range m {
      buf = appendField(buf, key, val)
   }
   return buf
}

func (m Message) MarshalJSON() ([]byte, error) {
   mes := map[protowire.Number]interface{}(m)
   return json.Marshal(mes)
}

func (m *Message) UnmarshalJSON(buf []byte) error {
   var raw map[protowire.Number]json.RawMessage
   err := json.Unmarshal(buf, &raw)
   if err != nil {
      return err
   }
   for key, val := range raw {
      any, err := unmarshal(val)
      if err != nil {
         return err
      }
      (*m)[key] = any
   }
   return nil
}
