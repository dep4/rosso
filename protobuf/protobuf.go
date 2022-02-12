// Protocol Buffers
package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
   "strings"
)

const (
   bytesType = "bytes"
   fixed64Type = "fixed64"
   messageType = "message"
   stringType = "string"
   varintType = "varint"
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
      for _, elem := range val {
         buf = appendField(buf, num, elem)
      }
   case []string:
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

type Message map[Tag]interface{}

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
   for tag, val := range m {
      buf = appendField(buf, tag.Number, val)
   }
   return buf
}

type Tag struct {
   protowire.Number
   Name string
}

// encoding/json
func (t Tag) MarshalText() ([]byte, error) {
   var buf []byte
   buf = strconv.AppendInt(buf, int64(t.Number), 10)
   buf = append(buf, ' ')
   buf = append(buf, t.Name...)
   return buf, nil
}

type nilMap struct {
   value string
}

func (n nilMap) Error() string {
   return strconv.Quote(n.value) + " assignment to entry in nil map"
}
