package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case bool:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, protowire.EncodeBool(val))
   case uint32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []byte:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val)
   case []interface{}:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
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
      mes := Unmarshal(buf)
      if mes != nil {
         return mes, vLen
      }
      if isBinary(buf) {
         return buf, vLen
      }
      return string(buf), vLen
   case protowire.StartGroupType:
      buf, vLen := protowire.ConsumeGroup(num, buf)
      mes := Unmarshal(buf)
      if mes != nil {
         return mes, vLen
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

type Message map[protowire.Number]interface{}

func Decode(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf), nil
}

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

func (m Message) Encode() io.Reader {
   buf := m.Marshal()
   return bytes.NewReader(buf)
}

func (m Message) Get(k protowire.Number) Message {
   val, ok := m[k].(Message)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetMessages(k protowire.Number) []Message {
   switch typ := m[k].(type) {
   case Message:
      return []Message{typ}
   case []interface{}:
      var mess []Message
      for _, val := range typ {
         mes, ok := val.(Message)
         if ok {
            mess = append(mess, mes)
         }
      }
      return mess
   default:
      return nil
   }
}

func (m Message) GetString(k protowire.Number) string {
   val, ok := m[k].(string)
   if ok {
      return val
   }
   return ""
}

func (m Message) GetUint64(k protowire.Number) uint64 {
   val, ok := m[k].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) Marshal() []byte {
   var buf []byte
   for key, val := range m {
      buf = appendField(buf, key, val)
   }
   return buf
}

func (m Message) Set(k protowire.Number, v interface{}) bool {
   if m == nil {
      return false
   }
   m[k] = v
   return true
}
