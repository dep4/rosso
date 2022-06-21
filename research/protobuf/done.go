package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "sort"
   "strconv"
)

func (m Message) consume_varint(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := add(m, num, val); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func add[T Encoder](mes Message, num Number, val T) error {
   in := mes[num]
   switch out := in.(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Slice[T]{out, val}
   case Slice[T]:
      mes[num] = append(out, val)
   default:
      return type_error{num, in, out}
   }
   return nil
}

type Slice[T Encoder] []T

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

func (s Slice[T]) encode(buf []byte, num Number) []byte {
   for _, encoder := range s {
      buf = encoder.encode(buf, num)
   }
   return buf
}

func (t type_error) Error() string {
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(t.Number), 10)
   b = append(b, " is "...)
   b = append(b, t.in.get_type()...)
   b = append(b, ", not "...)
   b = append(b, t.out.get_type()...)
   return string(b)
}

func (m Message) Marshal() []byte {
   var (
      nums []Number
      bufs []byte
   )
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      bufs = m[num].encode(bufs, num)
   }
   return bufs
}

type Number = protowire.Number

type type_error struct {
   Number
   in Encoder
   out Encoder
}

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

type Varint uint64

func (Varint) get_type() string { return "Varint" }

func (v Varint) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}

type Fixed32 uint32

func (Fixed32) get_type() string { return "Fixed32" }

func (f Fixed32) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

type Fixed64 uint64

func (Fixed64) get_type() string { return "Fixed64" }

func (f Fixed64) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

type Bytes []byte

func (Bytes) get_type() string { return "Bytes" }

func (b Bytes) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b)
}

type Message map[Number]Encoder

func (Message) get_type() string { return "Message" }

func (m Message) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

type Raw struct {
   Bytes []byte
   Message Message
}

func (Raw) get_type() string { return "Raw" }

func (r Raw) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, r.Bytes)
}
