package protobuf

import (
   "bytes"
   "encoding/base64"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Bytes struct {
   Raw Raw // Do not embed to keep MarshalText scoped to this field
   Message
}

func String(s string) Bytes {
   var dst Bytes
   dst.Raw = []byte(s)
   return dst
}

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Token

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

func (m Message) Add(num Number, val Message) {
   add(m, num, val)
}

func (m Message) AddString(num Number, val string) {
   add(m, num, String(val))
}

// Check Bytes for Unmarshaled Messages, check Message for manually constructed
// Messages.
func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Bytes:
      return value.Message
   case Message:
      return value
   }
   return nil
}

func (m Message) GetBytes(num Number) ([]byte, error) {
   src := m[num]
   dst, ok := src.(Bytes)
   if !ok {
      return nil, getError{src, num, dst}
   }
   return dst.Raw, nil
}

func (m Message) GetFixed64(num Number) (uint64, error) {
   src := m[num]
   dst, ok := src.(Fixed64)
   if !ok {
      return 0, getError{src, num, dst}
   }
   return uint64(dst), nil
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

func (m Message) GetString(num Number) (string, error) {
   src := m[num]
   dst, ok := src.(Bytes)
   if !ok {
      return "", getError{src, num, dst}
   }
   return string(dst.Raw), nil
}

func (m Message) GetVarint(num Number) (uint64, error) {
   src := m[num]
   dst, ok := src.(Varint)
   if !ok {
      return 0, getError{src, num, dst}
   }
   return uint64(dst), nil
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

type Raw []byte

func (r Raw) MarshalText() ([]byte, error) {
   if format.IsString(r) {
      return r, nil
   }
   buf := new(bytes.Buffer)
   base64.NewEncoder(base64.StdEncoding, buf).Write(r)
   return buf.Bytes(), nil
}

type Token interface {
   appendField([]byte, Number) []byte
}

type Tokens[T Token] []T

type Varint uint64
