package protobuf

import (
   "encoding/base64"
   "fmt"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strings"
)

func (m Message) Add(num Number, val Message) {
   add(m, num, val)
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
      buf = m[num].appendField(buf, num)
   }
   return buf
}

type Number = protowire.Number

type Token interface {
   appendField([]byte, Number) []byte
}

type Tokens[T Token] []T

type Varint uint64

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Token

type Bytes struct {
   Raw
   Message
}

type Raw []byte

func (r Raw) String() string {
   if format.IsString(r) {
      return string(r)
   }
   return base64.StdEncoding.EncodeToString(r)
}

type getError struct {
   src Token
   Number
   dst Token
}

func (g getError) Error() string {
   b := new(strings.Builder)
   fmt.Fprintf(b, "cannot unmarshal %T", g.src)
   fmt.Fprintf(b, " into field %v", g.Number)
   fmt.Fprintf(b, " of type %T", g.dst)
   return b.String()
}

func (m Message) GetFixed64(num Number) (uint64, error) {
   src := m[num]
   dst, ok := src.(Fixed64)
   if !ok {
      return 0, getError{src, num, dst}
   }
   return uint64(dst), nil
}

func (m Message) GetVarint(num Number) (uint64, error) {
   src := m[num]
   dst, ok := src.(Varint)
   if !ok {
      return 0, getError{src, num, dst}
   }
   return uint64(dst), nil
}

func (m Message) GetBytes(num Number) ([]byte, error) {
   src := m[num]
   dst, ok := src.(Bytes)
   if !ok {
      return nil, getError{src, num, dst}
   }
   return dst.Raw, nil
}

func add[T Token](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Tokens[T]{value, val}
   case Tokens[T]:
      mes[num] = append(value, val)
   }
}

func (m Message) Get(num Number) Message {
   dst, ok := m[num].(Bytes)
   if !ok {
      return nil
   }
   return dst.Message
}

func (m Message) AddBytes(num Number, val []byte) {
   dst := Bytes{Raw: val}
   add(m, num, dst)
}

func (f Fixed32) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

func (f Fixed64) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

func (t Tokens[T]) appendField(buf []byte, num Number) []byte {
   for _, tok := range t {
      buf = tok.appendField(buf, num)
   }
   return buf
}

func (v Varint) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}

func (b Bytes) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b.Raw)
}

func (m Message) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

func (m Message) GetMessages(num Number) []Message {
   var mes []Message
   switch value := m[num].(type) {
   case Bytes:
      return []Message{value.Message}
   case Tokens[Bytes]:
      for _, val := range value {
         mes = append(mes, val.Message)
      }
   }
   return mes
}

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, tLen := protowire.ConsumeTag(buf)
      err := protowire.ParseError(tLen)
      if err != nil {
         return nil, err
      }
      buf = buf[tLen:]
      var vLen int
      switch typ {
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(buf)
         add(mes, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(buf)
         add(mes, num, Fixed64(val))
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(buf)
         add(mes, num, Varint(val))
      case protowire.StartGroupType:
         var val []byte
         val, vLen = protowire.ConsumeGroup(num, buf)
         embed, err := Unmarshal(val)
         if err != nil {
            return nil, err
         }
         add(mes, num, embed)
      case protowire.BytesType:
         var val []byte
         val, vLen = protowire.ConsumeBytes(buf)
         embed, _ := Unmarshal(val)
         add(mes, num, Bytes{val, embed})
      }
      if err := protowire.ParseError(vLen); err != nil {
         return nil, err
      }
      buf = buf[vLen:]
   }
   return mes, nil
}
