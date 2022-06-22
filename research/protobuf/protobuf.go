package protobuf

import (
   "errors"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strconv"
)

type Message map[Number]Encoder

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, length := protowire.ConsumeTag(buf)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      buf = buf[length:]
      switch typ {
      case protowire.VarintType:
         buf, err = mes.consume_varint(num, buf)
      case protowire.Fixed64Type:
         buf, err = mes.consume_fixed64(num, buf)
      case protowire.Fixed32Type:
         buf, err = mes.consume_fixed32(num, buf)
      case protowire.BytesType:
         buf, err = mes.consume_raw(num, buf)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
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

func (m Message) consume_fixed32(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed32(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := Add(m, num, Fixed32(val)); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) consume_fixed64(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeFixed64(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := Add(m, num, Fixed64(val)); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) consume_raw(num Number, buf []byte) ([]byte, error) {
   var (
      length int
      val Raw
   )
   val.Bytes, length = protowire.ConsumeBytes(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if format.String(val.Bytes) {
      val.String = string(val.Bytes)
   }
   val.Message, _ = Unmarshal(val.Bytes)
   if err := Add(m, num, val); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) consume_varint(num Number, buf []byte) ([]byte, error) {
   val, length := protowire.ConsumeVarint(buf)
   err := protowire.ParseError(length)
   if err != nil {
      return nil, err
   }
   if err := Add(m, num, Varint(val)); err != nil {
      return nil, err
   }
   return buf[length:], nil
}

func (m Message) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

func (Message) get_type() string { return "Message" }

func Add[T Encoded](mes Message, num Number, val T) error {
   in := mes[num]
   switch out := in.(type) {
   case nil:
      mes[num] = Encoder(val)
   case T:
      mes[num] = Slice[T]{out, val}
   case Slice[T]:
      mes[num] = append(out, val)
   default:
      return type_error{num, in, out}
   }
   return nil
}

type Bytes []byte

func (b Bytes) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b)
}

func (Bytes) get_type() string { return "Bytes" }

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

type Fixed32 uint32

func (f Fixed32) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

func (Fixed32) get_type() string { return "Fixed32" }

type Fixed64 uint64

func (f Fixed64) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

func (Fixed64) get_type() string { return "Fixed64" }

type Number = protowire.Number

type Raw struct {
   Bytes []byte
   String string
   Message map[Number]Encoder
}

func (r Raw) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, r.Bytes)
}

func (Raw) get_type() string { return "Raw" }

type String string

func (s String) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendString(buf, string(s))
}

func (String) get_type() string { return "String" }

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

type Encoded interface {
   Fixed32 | Fixed64 | Message | Raw | String | Varint
}

type Slice[T Encoded] []T

func (s Slice[T]) encode(buf []byte, num Number) []byte {
   for _, value := range s {
      buf = Encoder(value).encode(buf, num)
   }
   return buf
}

func (Slice[T]) get_type() string {
   var value T
   return "[]" + Encoder(value).get_type()
}
