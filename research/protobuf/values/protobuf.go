package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strconv"
)

type Encoder interface {
   encode([]byte, Number) []byte
   valueType() string
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

func (m Message) consumeBytes(num Number, buf []byte) ([]byte, error) {
   vals, err := get[SliceBytes](m, num)
   if err != nil {
      return nil, err
   }
   var (
      val Bytes
      vLen int
   )
   val.Raw, vLen = protowire.ConsumeBytes(buf)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   val.Message, err = Unmarshal(val.Raw)
   if err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[vLen:], nil
}

type Number = protowire.Number

func (SliceBytes) valueType() string { return "SliceBytes" }

func (s SliceBytes) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Raw)
   }
   return buf
}

type SliceFixed32 []uint32

func (s SliceFixed32) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   }
   return buf
}

func (SliceFixed32) valueType() string { return "SliceFixed32" }

type SliceFixed64 []uint64

func (s SliceFixed64) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
      buf = protowire.AppendFixed64(buf, val)
   }
   return buf
}

func (SliceFixed64) valueType() string { return "SliceFixed64" }

func (s SliceMessage) encode(buf []byte, num Number) []byte {
   for _, mes := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, mes.Marshal())
   }
   return buf
}

func (SliceMessage) valueType() string { return "SliceMessage" }

type SliceVarint []uint64

func (s SliceVarint) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   }
   return buf
}

func (SliceVarint) valueType() string { return "SliceVarint" }

////////////////////////////

type SliceMessage []Message

func (t typeError) Error() string {
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(t.Number), 10)
   b = append(b, " is "...)
   b = append(b, t.in.valueType()...)
   b = append(b, ", not "...)
   b = append(b, t.out.valueType()...)
   return string(b)
}

type typeError struct {
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
   return out, typeError{num, in, out}
}

func (m Message) consumeVarint(num Number, buf []byte) ([]byte, error) {
   vals, err := get[SliceVarint](m, num)
   if err != nil {
      return nil, err
   }
   val, vLen := protowire.ConsumeVarint(buf)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[vLen:], nil
}

func (m Message) consumeFixed64(num Number, buf []byte) ([]byte, error) {
   vals, err := get[SliceFixed64](m, num)
   if err != nil {
      return nil, err
   }
   val, vLen := protowire.ConsumeFixed64(buf)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[vLen:], nil
}

func (m Message) consumeFixed32(num Number, buf []byte) ([]byte, error) {
   vals, err := get[SliceFixed32](m, num)
   if err != nil {
      return nil, err
   }
   val, vLen := protowire.ConsumeFixed32(buf)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return buf[vLen:], nil
}

type SliceBytes []Bytes

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
      num, typ, tLen := protowire.ConsumeTag(buf)
      err := protowire.ParseError(tLen)
      if err != nil {
         return nil, err
      }
      buf = buf[tLen:]
      switch typ {
      case protowire.VarintType:
         buf, err = mes.consumeVarint(num, buf)
      case protowire.Fixed64Type:
         buf, err = mes.consumeFixed64(num, buf)
      case protowire.Fixed32Type:
         buf, err = mes.consumeFixed32(num, buf)
      case protowire.BytesType:
         buf, err = mes.consumeBytes(num, buf)
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}
