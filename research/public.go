package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
)

type Message map[Tag]interface{}

func Decode(r io.Reader) (Message, error) {
   b, err := io.ReadAll(r)
   if err != nil {
      return nil, err
   }
   return Unmarshal(b)
}

func Unmarshal(b []byte) (Message, error) {
   mes := make(Message)
   for len(b) >= 1 {
      num, typ, fLen := protowire.ConsumeField(b)
      err := protowire.ParseError(fLen)
      if err != nil {
         return nil, err
      }
      _, _, tLen := protowire.ConsumeTag(b[:fLen])
      if err := protowire.ParseError(tLen); err != nil {
         return nil, err
      }
      value := b[tLen:fLen]
      switch typ {
      case protowire.BytesType:
         err = mes.consumeBytes(num, value)
      case protowire.Fixed64Type:
         err = mes.consumeFixed64(num, value)
      case protowire.VarintType:
         err = mes.consumeVarint(num, value)
      }
      if err != nil {
         return nil, err
      }
      b = b[fLen:]
   }
   return mes, nil
}

func (m Message) Add(num protowire.Number, v Message) {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = v
   case Message:
      m[tag] = []Message{value, v}
   case []Message:
      m[tag] = append(value, v)
   }
}

func (m Message) Get(num protowire.Number) Message {
   value, ok := m[Tag{num, messageType}].(Message)
   if ok {
      return value
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number) []byte {
   value, ok := m[Tag{num, bytesType}].([]byte)
   if ok {
      return value
   }
   return nil
}

func (m Message) GetFixed64(num protowire.Number) uint64 {
   value, ok := m[Tag{num, fixed64Type}].(uint64)
   if ok {
      return value
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
   value, ok := m[Tag{num, bytesType}].(string)
   if ok {
      return value
   }
   return ""
}

func (m Message) GetVarint(num protowire.Number) uint64 {
   value, ok := m[Tag{num, varintType}].(uint64)
   if ok {
      return value
   }
   return 0
}

func (m Message) Marshal() []byte {
   var buf []byte
   for tag, value := range m {
      buf = appendField(buf, tag.Number, value)
   }
   return buf
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
