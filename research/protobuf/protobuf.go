package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "sort"
   "strconv"
)

type Message map[Number]Encoder

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

func (m Message) Add(num Number, val Encoder) error {
   return add(m, num, val)
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

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

type Number = protowire.Number

type Slice[T Encoder] []T

func (s Slice[T]) encode(buf []byte, num Number) []byte {
   for _, encoder := range s {
      buf = encoder.encode(buf, num)
   }
   return buf
}

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

type Varint uint64

func (v Varint) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}

func (Varint) get_type() string { return "Varint" }

type type_error struct {
   Number
   in Encoder
   out Encoder
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
