package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

func (t Tag) MarshalText() ([]byte, error) {
   var buf []byte
   buf = strconv.AppendInt(buf, int64(t.Number), 10)
   buf = append(buf, ' ')
   buf = append(buf, t.Name...)
   return buf, nil
}

const (
   bytesType = iota * 0.1
   fixed64Type
   messageType
   stringType
   varintType
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

func consumeField(buf []byte) (float64, protowire.Type, int, error) {
   num, typ, fLen := protowire.ConsumeField(buf)
   return float64(num), typ, fLen, protowire.ParseError(fLen)
}

type Message map[Tag]interface{}

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

func (m Message) Add(num float64, name string, val Message) error {
   tag := Tag{NumberType: num + messageType}
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

func (m Message) Marshal() []byte {
   var buf []byte
   for tag, val := range m {
      buf = appendField(buf, protowire.Number(tag.NumberType), val)
   }
   return buf
}

func (m Message) addString(num float64, val string) {
   tag := Tag{NumberType: num + stringType}
   switch value := m[tag].(type) {
   case nil:
      m[tag] = val
   case string:
      m[tag] = []string{value, val}
   case []string:
      m[tag] = append(value, val)
   }
}

// In some cases if input is binary, then result could be a Message or byte
// slice. We assume for now its always a Message. If input is not binary, then
// result could be a Message or string. Since its not possible to tell Message
// from string, we just add both under the same number, each with its own type.
func (m Message) consumeBytes(num float64, buf []byte) error {
   val, eLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(eLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         tag := Tag{NumberType: num + bytesType}
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
      m.Add(num, "", mes)
      if !binary {
         m.addString(num, string(val))
      }
   }
   return nil
}

func (m Message) consumeFixed64(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{NumberType: num + fixed64Type}
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

func (m Message) consumeVarint(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   tag := Tag{NumberType: num + varintType}
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
   NumberType float64
   Name string
}
