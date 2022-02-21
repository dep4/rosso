package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
   "strings"
)

// We cannot include the name in the key. When you Unmarshal, the name will be
// empty. If you then try to Get with a name, it will fail. Max valid number is
// 536,870,911, so better to use float64:
// stackoverflow.com/questions/3793838
type Message map[Number]interface{}

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

func (m Message) Add(num Number, name string, val Message) error {
   num += messageType
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

func (m Message) Get(num Number, name string) Message {
   val, ok := m[num + messageType].(Message)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetBytes(num Number, name string) []byte {
   val, ok := m[num + bytesType].([]byte)
   if ok {
      return val
   }
   return nil
}

func (m Message) GetFixed64(num Number, name string) uint64 {
   val, ok := m[num + fixed64Type].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) GetMessages(num Number, name string) []Message {
   switch value := m[num + messageType].(type) {
   case []Message:
      return value
   case Message:
      return []Message{value}
   }
   return nil
}

func (m Message) GetString(num Number, name string) string {
   val, ok := m[num + bytesType].(string)
   if ok {
      return val
   }
   return ""
}

func (m Message) GetVarint(num Number, name string) uint64 {
   val, ok := m[num + varintType].(uint64)
   if ok {
      return val
   }
   return 0
}

func (m Message) GoString() string {
   buf := new(strings.Builder)
   buf.WriteString("protobuf.Message{")
   first := true
   for tag, val := range m {
      if first {
         first = false
      } else {
         buf.WriteString(",\n")
      }
      fmt.Fprintf(buf, "%#v:", tag)
      num, ok := val.(uint64)
      if ok {
         fmt.Fprintf(buf, "uint64(%v)", num)
      } else {
         fmt.Fprintf(buf, "%#v", val)
      }
   }
   buf.WriteByte('}')
   return buf.String()
}

func (m Message) Marshal() []byte {
   var buf []byte
   for num, val := range m {
      buf = appendField(buf, protowire.Number(num), val)
   }
   return buf
}

type Number float64

func Tag(num Number, name string) Number {
   return num
}

func (n Number) MarshalText() ([]byte, error) {
   f := float64(n)
   return strconv.AppendFloat(nil, f, 'f', -1, 64), nil
}
