package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
)

const (
   MessageType = 0
   BytesType = 0.1
   VarintType = 0.2
   Fixed64Type = 0.3
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
      num, typ, fLen, err := consumeField(buf)
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

// Add value using MessageType, given number and empty Name. To add value with
// a Name, use the map directly.
func (m Message) Add(num float64, val Message) error {
   tag := Tag{NumberType: num + MessageType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case Message:
      m[tag] = []Message{value, val}
   case []Message:
      m[tag] = append(value, val)
   }
   return nil
}

// Return value using MessageType, given number and empty Name. To return value
// with a Name, use the map directly.
func (m Message) Get(num float64) Message {
   tag := Tag{NumberType: num + MessageType}
   val, ok := m[tag].(Message)
   if ok {
      return val
   }
   return nil
}

// Return value using BytesType, given number and empty Name. To return value
// with a Name, use the map directly.
func (m Message) GetBytes(num float64) []byte {
   tag := Tag{NumberType: num + BytesType}
   val, ok := m[tag].([]byte)
   if ok {
      return val
   }
   return nil
}

// Return value using Fixed64Type, given number and empty Name. To return value
// with a Name, use the map directly.
func (m Message) GetFixed64(num float64) uint64 {
   tag := Tag{NumberType: num + Fixed64Type}
   val, ok := m[tag].(uint64)
   if ok {
      return val
   }
   return 0
}

// Return value using MessageType, given number and empty Name. To return value
// with a Name, use the map directly.
func (m Message) GetMessages(num float64) []Message {
   tag := Tag{NumberType: num + MessageType}
   val, ok := m[tag].([]Message)
   if ok {
      return val
   }
   return nil
}

// Return value using BytesType, given number and empty Name. To return value
// with a Name, use the map directly.
func (m Message) GetString(num float64) string {
   tag := Tag{NumberType: num + BytesType}
   val, ok := m[tag].(string)
   if ok {
      return val
   }
   return ""
}

// Return value using VarintType, given number and empty Name. To return value
// with a Name, use the map directly.
func (m Message) GetVarint(num float64) uint64 {
   tag := Tag{NumberType: num + VarintType}
   val, ok := m[tag].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) Marshal() []byte {
   var buf []byte
   for tag, val := range m {
      buf = appendField(buf, protowire.Number(tag.NumberType), val)
   }
   return buf
}

type Tag struct {
   NumberType float64
   Name string
}

func (t Tag) MarshalText() ([]byte, error) {
   return strconv.AppendFloat(nil, t.NumberType, 'f', -1, 64), nil
}
