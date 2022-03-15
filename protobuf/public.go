package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
)

type Message map[Tag]interface{}

func Decode(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

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
      value := buf[tLen:fLen]
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
      buf = buf[fLen:]
   }
   return mes, nil
}

func (m Message) Add(num protowire.Number, s string, v Message) {
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

func (m Message) Get(num protowire.Number, s string) Message {
   tag := Tag{num, messageType}
   value, ok := m[tag].(Message)
   if ok {
      return value
   }
   return nil
}

func (m Message) GetBytes(num protowire.Number, s string) []byte {
   tag := Tag{num, bytesType}
   value, ok := m[tag].([]byte)
   if ok {
      return value
   }
   return nil
}

func (m Message) GetFixed64(num protowire.Number, s string) uint64 {
   tag := Tag{num, fixed64Type}
   value, ok := m[tag].(uint64)
   if ok {
      return value
   }
   return 0
}

func (m Message) GetMessages(num protowire.Number, s string) []Message {
   tag := Tag{num, messageType}
   switch value := m[tag].(type) {
   case []Message:
      return value
   case Message:
      return []Message{value}
   }
   return nil
}

func (m Message) GetString(num protowire.Number, s string) string {
   tag := Tag{num, bytesType}
   value, ok := m[tag].(string)
   if ok {
      return value
   }
   return ""
}

func (m Message) GetVarint(num protowire.Number, s string) uint64 {
   tag := Tag{num, varintType}
   value, ok := m[tag].(uint64)
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

func MessageTag(num protowire.Number, s string) Tag {
   return Tag{num, messageType}
}

func NewTag(num protowire.Number, s string) Tag {
   return Tag{Number: num}
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
