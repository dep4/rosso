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

type SliceVarint []uint64

type SliceFixed64 []uint64

type SliceFixed32 []uint32

type Number = protowire.Number

func (s SliceVarint) encode(b []byte, num Number) []byte {
   for _, val := range s {
      b = protowire.AppendTag(b, num, protowire.VarintType)
      b = protowire.AppendVarint(b, val)
   }
   return b
}

func (s SliceFixed64) encode(b []byte, num Number) []byte {
   for _, val := range s {
      b = protowire.AppendTag(b, num, protowire.Fixed64Type)
      b = protowire.AppendFixed64(b, val)
   }
   return b
}

func (s SliceFixed32) encode(b []byte, num Number) []byte {
   for _, val := range s {
      b = protowire.AppendTag(b, num, protowire.Fixed32Type)
      b = protowire.AppendFixed32(b, val)
   }
   return b
}

func (s SliceBytes) encode(b []byte, num Number) []byte {
   for _, val := range s {
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendBytes(b, val.Raw)
   }
   return b
}

func (s SliceMessage) encode(b []byte, num Number) []byte {
   for _, mes := range s {
      b = protowire.AppendTag(b, num, protowire.BytesType)
      b = protowire.AppendBytes(b, mes.MarshalBinary())
   }
   return b
}

func (m Message) MarshalBinary() []byte {
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

func (SliceBytes) valueType() string { return "SliceBytes" }

func (SliceFixed32) valueType() string { return "SliceFixed32" }

func (SliceFixed64) valueType() string { return "SliceFixed64" }

func (SliceVarint) valueType() string { return "SliceVarint" }

func (SliceMessage) valueType() string { return "SliceMessage" }

func Unmarshal(b []byte) (Message, error) {
   if len(b) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(b) >= 1 {
      num, typ, tLen := protowire.ConsumeTag(b)
      err := protowire.ParseError(tLen)
      if err != nil {
         return nil, err
      }
      b = b[tLen:]
      switch typ {
      case protowire.VarintType:
         b, err = mes.consumeVarint(num, b)
      case protowire.Fixed64Type:
         b, err = mes.consumeFixed64(num, b)
      case protowire.Fixed32Type:
         b, err = mes.consumeFixed32(num, b)
      case protowire.BytesType:
         b, err = mes.consumeBytes(num, b)
      }
      if err != nil {
         return nil, err
      }
   }
   return mes, nil
}

type Message map[Number]Encoder

type SliceMessage []Message

type Bytes struct {
   Raw []byte
   Message Message
}

type SliceBytes []Bytes

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

func (m Message) consumeVarint(num Number, b []byte) ([]byte, error) {
   vals, err := get[SliceVarint](m, num)
   if err != nil {
      return nil, err
   }
   val, vLen := protowire.ConsumeVarint(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return b[vLen:], nil
}

func (m Message) consumeFixed64(num Number, b []byte) ([]byte, error) {
   vals, err := get[SliceFixed64](m, num)
   if err != nil {
      return nil, err
   }
   val, vLen := protowire.ConsumeFixed64(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return b[vLen:], nil
}

func (m Message) consumeFixed32(num Number, b []byte) ([]byte, error) {
   vals, err := get[SliceFixed32](m, num)
   if err != nil {
      return nil, err
   }
   val, vLen := protowire.ConsumeFixed32(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   m[num] = append(vals, val)
   return b[vLen:], nil
}

/////////////////////////////////////////////////////////////////////////////

func (m Message) consumeBytes(num Number, b []byte) ([]byte, error) {
   /*
   var val Bytes
   val.Message = make(Message)
   val.Raw, vLen = protowire.ConsumeBytes(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   err := val.Message.UnmarshalBinary(val.Raw)
   if err != nil {
      val.Message = nil
   }
   mes[num] = append(mes[num], val)
   */
   val, vLen := protowire.ConsumeBytes(b)
   if err := protowire.ParseError(vLen); err != nil {
      return nil, err
   }
   if in := m[num]; in == nil {
      m[num] = SliceBytes{Bytes{Raw: val}}
   } else if out, ok := in.(SliceBytes); ok {
      m[num] = append(out, Bytes{Raw: val})
   } else {
      return nil, typeError{num, in, out}
   }
   return b[vLen:], nil
}
