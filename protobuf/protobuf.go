package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case uint32:
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case []interface{}:
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   }
   return buf
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
   return Unmarshal(buf)
}

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen, err := consumeField(buf)
      if err != nil {
         return nil, err
      }
      tLen, err := consumeTag(buf[:fLen])
      if err != nil {
         return nil, err
      }
      val, err := consume(num, typ, buf[tLen:fLen])
      if err != nil {
         return nil, err
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
   return mes, nil
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
   case []Message:
      return typ
   case Message:
      return []Message{typ}
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
