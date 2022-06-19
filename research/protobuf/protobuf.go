package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strconv"
)

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

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

func (m Message) consume_bytes(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Slice_Bytes](m, num)
   if err != nil {
      return nil, err
   }
   var (
      val Bytes
      v_len int
   )
   val.Raw, v_len = protowire.ConsumeBytes(buf)
   if err := protowire.ParseError(v_len); err != nil {
      return nil, err
   }
   val.Message, err = Unmarshal(val.Raw)
   if err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[v_len:], nil
}

type Number = protowire.Number

func (Slice_Bytes) get_type() string { return "Slice_Bytes" }

func (s Slice_Bytes) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Raw)
   }
   return buf
}

type Slice_Fixed32 []uint32

func (s Slice_Fixed32) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   }
   return buf
}

func (Slice_Fixed32) get_type() string { return "Slice_Fixed32" }

type Slice_Fixed64 []uint64

func (s Slice_Fixed64) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
      buf = protowire.AppendFixed64(buf, val)
   }
   return buf
}

func (Slice_Fixed64) get_type() string { return "Slice_Fixed64" }

func (s Slice_Message) encode(buf []byte, num Number) []byte {
   for _, mes := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, mes.Marshal())
   }
   return buf
}

func (Slice_Message) get_type() string { return "Slice_Message" }

type Slice_Varint []uint64

func (s Slice_Varint) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   }
   return buf
}

func (Slice_Varint) get_type() string { return "Slice_Varint" }

type Slice_Message []Message

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

type type_error struct {
   Number
   in Encoder
   out Encoder
}

func get[T Encoder](mes Message, num Number) (T, error) {
   in := mes[num]
   out, ok := in.(T)
   if in == nil {
      return out, nil
   }
   if ok {
      return out, nil
   }
   return out, type_error{num, in, out}
}

func (m Message) consume_varint(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Slice_Varint](m, num)
   if err != nil {
      return nil, err
   }
   val, v_len := protowire.ConsumeVarint(buf)
   if err := protowire.ParseError(v_len); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[v_len:], nil
}

func (m Message) consume_fixed64(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Slice_Fixed64](m, num)
   if err != nil {
      return nil, err
   }
   val, v_len := protowire.ConsumeFixed64(buf)
   if err := protowire.ParseError(v_len); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[v_len:], nil
}

func (m Message) consume_fixed32(num Number, buf []byte) ([]byte, error) {
   vals, err := get[Slice_Fixed32](m, num)
   if err != nil {
      return nil, err
   }
   val, v_len := protowire.ConsumeFixed32(buf)
   if err := protowire.ParseError(v_len); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[v_len:], nil
}

type Slice_Bytes []Bytes

type Bytes struct {
   Raw []byte
   Message Message
}

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, t_len := protowire.ConsumeTag(buf)
      err := protowire.ParseError(t_len)
      if err != nil {
         return nil, err
      }
      buf = buf[t_len:]
      switch typ {
      case protowire.VarintType:
         buf, err = mes.consume_varint(num, buf)
      case protowire.Fixed64Type:
         buf, err = mes.consume_fixed64(num, buf)
      case protowire.Fixed32Type:
         buf, err = mes.consume_fixed32(num, buf)
      case protowire.BytesType:
         buf, err = mes.consume_bytes(num, buf)
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}
