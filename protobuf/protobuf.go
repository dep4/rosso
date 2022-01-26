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
   fixed32Type = "fixed32"
   groupType = "group"
   messageType = "message"
   stringType = "string"
   varintType = "varint"
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
      case protowire.Fixed32Type:
         err = mes.consumeFixed32(num, val)
      case protowire.StartGroupType:
         err = mes.consumeGroup(num, val)
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
   str := new(strings.Builder)
   str.WriteString("protobuf.Message{")
   first := true
   for tag, val := range m {
      if first {
         first = false
      } else {
         str.WriteString(",\n")
      }
      fmt.Fprintf(str, "%#v:", tag)
      switch typ := val.(type) {
      case uint32:
         fmt.Fprintf(str, "uint32(%v)", typ)
      case uint64:
         fmt.Fprintf(str, "uint64(%v)", typ)
      default:
         fmt.Fprintf(str, "%#v", val)
      }
   }
   str.WriteByte('}')
   return str.String()
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
