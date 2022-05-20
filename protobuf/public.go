package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Bytes []byte

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
      case protowire.BytesType:
         var val []byte
         val, vLen = protowire.ConsumeBytes(buf)
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
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(buf)
         add(mes, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(buf)
         add(mes, num, Fixed64(val))
      case protowire.StartGroupType:
         var val []byte
         val, vLen = protowire.ConsumeGroup(num, buf)
         embed, err := Unmarshal(val)
         if err != nil {
            return nil, err
         }
         add(mes, num, embed)
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(buf)
         add(mes, num, Varint(val))
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
