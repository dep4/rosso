package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

type Bytes []byte

type Message map[protowire.Number]interface{}

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

func (m Message) Get(k protowire.Number) Message {
   val, ok := m[k].(Message)
   if ok {
      return val
   }
   return nil
}

// mimesniff.spec.whatwg.org#binary-data-byte
func (b Bytes) String() string {
   for _, c := range b {
      switch {
      case c <= 0x08,
      c == 0x0B,
      0x0E <= c && c <= 0x1A,
      0x1C <= c && c <= 0x1F:
         return fmt.Sprint([]byte(b))
      }
   }
   return string(b)
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
      return Bytes(buf), vLen
   case protowire.StartGroupType:
      buf, vLen := protowire.ConsumeGroup(num, buf)
      mes := Unmarshal(buf)
      if mes != nil {
         return mes, vLen
      }
      return Bytes(buf), vLen
   }
   return nil, 0
}

func (m Message) GetBytes(k protowire.Number) Bytes {
   val, ok := m[k].(Bytes)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetMessages(k protowire.Number) []Message {
   switch typ := m[k].(type) {
   case []Message:
      return typ
   case Message:
      return []Message{typ}
   default:
      return nil
   }
}

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case uint32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case Bytes:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, []byte(val))
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case []uint32:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []uint64:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []Bytes:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []Message:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   }
   return buf
}
