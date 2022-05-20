package protobuf

import (
   "fmt"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strings"
)

type Bytes []byte

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Token

func Decode(in io.Reader) (Message, error) {
   buf, err := io.ReadAll(in)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

func (m Message) Add(num Number, val Message) {
   add(m, num, val)
}

func (m Message) AddString(num Number, val String) {
   add(m, num, val)
}

func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Message:
      return value
   case String:
      return m.Get(-num)
   }
   return nil
}

func (m Message) GetBytes(num Number) (Bytes, error) {
   return get[Bytes](m, num)
}

func (m Message) GetFixed64(num Number) (Fixed64, error) {
   return get[Fixed64](m, num)
}

func (m Message) GetMessages(num Number) []Message {
   switch value := m[num].(type) {
   case Tokens[Message]:
      return value
   case Message:
      return []Message{value}
   }
   return nil
}

func (m Message) GetString(num Number) (String, error) {
   return get[String](m, num)
}

func (m Message) GetVarint(num Number) (Varint, error) {
   return get[Varint](m, num)
}

func (m Message) GoString() string {
   buf := new(strings.Builder)
   buf.WriteString("protobuf.Message{\n")
   for num, tok := range m {
      fmt.Fprint(buf, num, ":")
      switch tok.(type) {
      case Fixed32:
         fmt.Fprintf(buf, "protobuf.Fixed32(%v)", tok)
      case Fixed64:
         fmt.Fprintf(buf, "protobuf.Fixed64(%v)", tok)
      case String:
         fmt.Fprintf(buf, "protobuf.String(%q)", tok)
      case Varint:
         fmt.Fprintf(buf, "protobuf.Varint(%v)", tok)
      default:
         fmt.Fprintf(buf, "%#v", tok)
      }
      buf.WriteString(",\n")
   }
   buf.WriteByte('}')
   return buf.String()
}

func (m Message) Marshal() []byte {
   var (
      buf []byte
      nums []Number
   )
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      if num >= protowire.MinValidNumber {
         buf = m[num].appendField(buf, num)
      }
   }
   return buf
}

type Number = protowire.Number

type String string

type Token interface {
   appendField([]byte, Number) []byte
}

type Tokens[T Token] []T

type Varint uint64

func Unmarshal(b []byte) (Message, error) {
   mes := make(Message)
   for len(b) >= 1 {
      num, typ, n := protowire.ConsumeTag(b)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      b = b[n:]
      m := protowire.ConsumeFieldValue(num, typ, b)
      if err := protowire.ParseError(m); err != nil {
         return nil, err
      }
      switch typ {
      case protowire.BytesType:
         val, vLen := protowire.ConsumeBytes(b)
         err := protowire.ParseError(vLen)
         if err != nil {
            return nil, err
         }
         if len(val) == 0 {
            add(mes, num, String(""))
         } else {
            embed, err := Unmarshal(val)
            if err != nil {
               if format.IsString(val) {
                  add(mes, num, String(val))
               } else {
                  add(mes, num, Bytes(val))
               }
            } else if format.IsString(val) {
               add(mes, num, String(val))
               add(mes, -num, embed)
            } else {
               add(mes, num, embed)
            }
         }
      case protowire.Fixed32Type:
         val, vLen := protowire.ConsumeFixed32(b)
         err := protowire.ParseError(vLen)
         if err != nil {
            return nil, err
         }
         add(mes, num, Fixed32(val))
      case protowire.Fixed64Type:
         val, vLen := protowire.ConsumeFixed64(b)
         err := protowire.ParseError(vLen)
         if err != nil {
            return nil, err
         }
         add(mes, num, Fixed64(val))
      case protowire.VarintType:
         val, vLen := protowire.ConsumeVarint(b)
         err := protowire.ParseError(vLen)
         if err != nil {
            return nil, err
         }
         add(mes, num, Varint(val))
      }
      b = b[m:]
   }
   return mes, nil
}
