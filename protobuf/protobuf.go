package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

const (
   MessageType = 0
   BytesType = 0.1
   VarintType = 0.2
   Fixed64Type = 0.3
)

type Message map[float64]interface{}

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

// Add value using MessageType and given number.
func (m Message) Add(num float64, val Message) error {
   num += MessageType
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case Message:
      m[num] = []Message{value, val}
   case []Message:
      m[num] = append(value, val)
   }
   return nil
}

// Return value using MessageType and given number.
func (m Message) Get(num float64) Message {
   val, ok := m[num + MessageType].(Message)
   if ok {
      return val
   }
   return nil
}

// Return value using BytesType and given number.
func (m Message) GetBytes(num float64) []byte {
   val, ok := m[num + BytesType].([]byte)
   if ok {
      return val
   }
   return nil
}

// Return value using Fixed64Type and given number.
func (m Message) GetFixed64(num float64) uint64 {
   val, ok := m[num + Fixed64Type].(uint64)
   if ok {
      return val
   }
   return 0
}

// Return value using MessageType and given number.
func (m Message) GetMessages(num float64) []Message {
   val, ok := m[num + MessageType].([]Message)
   if ok {
      return val
   }
   return nil
}

// Return value using BytesType and given number.
func (m Message) GetString(num float64) string {
   val, ok := m[num + BytesType].(string)
   if ok {
      return val
   }
   return ""
}

// Return value using VarintType and given number.
func (m Message) GetVarint(num float64) uint64 {
   val, ok := m[num + VarintType].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, val := range m {
      buf = appendField(buf, protowire.Number(num), val)
   }
   return buf
}

func (m Message) addString(num float64, val string) {
   num += BytesType
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case string:
      m[num] = []string{value, val}
   case []string:
      m[num] = append(value, val)
   }
}

func (m Message) consumeFixed64(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   num += Fixed64Type
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case uint64:
      m[num] = []uint64{value, val}
   case []uint64:
      m[num] = append(value, val)
   }
   return nil
}

func (m Message) consumeVarint(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   num += VarintType
   switch value := m[num].(type) {
   case nil:
      m[num] = val
   case uint64:
      m[num] = []uint64{value, val}
   case []uint64:
      m[num] = append(value, val)
   }
   return nil
}

// In some cases if input is binary, then result could be a Message or byte
// slice. We assume for now its always a Message. If input is not binary, then
// result could be a Message or string. Since its not possible to tell Message
// from string, we just add both under the same number, each with its own type.
func (m Message) consumeBytes(num float64, buf []byte) error {
   val, vLen := protowire.ConsumeBytes(buf)
   err := protowire.ParseError(vLen)
   if err != nil {
      return err
   }
   binary := format.IsBinary(val)
   mes, err := Unmarshal(val)
   if err != nil {
      if binary {
         num += BytesType
         switch value := m[num].(type) {
         case nil:
            m[num] = val
         case []byte:
            m[num] = [][]byte{value, val}
         case [][]byte:
            m[num] = append(value, val)
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
