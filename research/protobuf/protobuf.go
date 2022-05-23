package protobuf

import (
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Token

type Number = protowire.Number

type Token interface {
   appendField([]byte, Number) []byte
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

type Varint uint64

type Tokens[T Token] []T

func (f Fixed32) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(in, uint32(f))
}

func (f Fixed64) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(in, uint64(f))
}

func (m Message) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendBytes(in, m.Marshal())
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

func (t Tokens[T]) appendField(in []byte, num Number) []byte {
   for _, tok := range t {
      in = tok.appendField(in, num)
   }
   return in
}

func (v Varint) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.VarintType)
   return protowire.AppendVarint(in, uint64(v))
}

func (s String) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendString(in, s.Raw)
}

func (b Bytes) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendBytes(in, b.Raw)
}

////////////////////////////////////////////////////

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
         embed, err := Unmarshal(val)
         if err != nil {
            if format.IsString(val) {
               add(mes, num, String{string(val), nil})
            } else {
               add(mes, num, Bytes{val, nil})
            }
         } else if format.IsString(val) {
            add(mes, num, String{string(val), embed})
         } else {
            add(mes, num, Bytes{val, embed})
         }
      }
      if err := protowire.ParseError(vLen); err != nil {
         return nil, err
      }
      buf = buf[vLen:]
   }
   return mes, nil
}

type String struct {
   Raw string
   Message
}

type Bytes struct {
   Raw []byte
   Message
}
