package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

const (
   bytesType = 2
   fixed64Type = 1
   messageType = 6
   varintType = 0
)

func appendField(buf []byte, num protowire.Number, val interface{}) []byte {
   switch val := val.(type) {
   case uint64:
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   case string:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendString(buf, val)
   case []byte:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val)
   case Message:
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Marshal())
   case []uint64:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []string:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   case []Message:
      for _, value := range val {
         buf = appendField(buf, num, value)
      }
   }
   return buf
}

type Message map[Tag]interface{}

func Unmarshal(buf []byte) (Message, error) {
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, fLen := protowire.ConsumeField(buf)
      err := protowire.ParseError(fLen)
      if err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(buf[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      val := buf[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         err = mes.consumeBytes(num, val)
      case protowire.Fixed64Type:
         err = mes.consumeFixed64(num, val)
      case protowire.VarintType:
         err = mes.consumeVarint(num, val)
      }
      if err != nil {
         return nil, err
      }
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Add(num protowire.Number, val Message) {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{value, val}
   case []Message:
      m[tag] = append(value, val)
   }
}

func (m Message) Get(num protowire.Number) Message {
   val, ok := m[Tag{num, messageType}].(Message)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number) []byte {
   val, ok := m[Tag{num, bytesType}].([]byte)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetFixed64(num protowire.Number) uint64 {
   val, ok := m[Tag{num, fixed64Type}].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) GetMessages(num protowire.Number) []Message {
   switch value := m[Tag{num, messageType}].(type) {
   case []Message:
      return value
   case Message:
      return []Message{value}
   }
   return nil
}

func (m Message) GetString(num protowire.Number) string {
   val, ok := m[Tag{num, bytesType}].(string)
   if ok {
      return val
   }
   return ""
}

func (m Message) GetVarint(num protowire.Number) uint64 {
   val, ok := m[Tag{num, varintType}].(uint64)
   if ok {
      return val
   }
   return 0
}
func (m Message) Marshal() []byte {
   var buf []byte
   for tag, val := range m {
      buf = appendField(buf, tag.Number, val)
   }
   return buf
}

func (m Message) addString(num protowire.Number, val string) {
   tag := Tag{num, bytesType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{value, val}
   case []string:
      m[tag] = append(value, val)
   }
}

func (m Message) consumeBytes(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         tag := Tag{num, bytesType}
         switch value := m[tag].(type) {
         case nil:
            m[tag] = val
         case []byte:
            m[tag] = [][]byte{value, val}
         case [][]byte:
            m[tag] = append(value, val)
         }
      } else {
         m.addString(num, string(val))
      }
   } else {
      m.Add(num, mes)
      if !binary {
         m.addString(num, string(val))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{num, fixed64Type}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{value, val}
   case []uint64:
      m[tag] = append(value, val)
   }
   return nil
}

func (m Message) consumeVarint(num protowire.Number, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{num, varintType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case uint64:
      m[tag] = []uint64{value, val}
   case []uint64:
      m[tag] = append(value, val)
   }
   return nil
}

type Tag struct {
   protowire.Number
   protowire.Type
}

func (t Tag) MarshalText() ([]byte, error) {
   buf := strconv.AppendInt(nil, int64(t.Number), 10)
   switch t.Type {
   case bytesType:
      buf = append(buf, " bytes"...)
   case fixed64Type:
      buf = append(buf, " fixed64"...)
   case messageType:
      buf = append(buf, " message"...)
   case varintType:
      buf = append(buf, " varint"...)
   }
   return buf, nil
}
