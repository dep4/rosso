package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Encoder interface {
   encode([]byte, Number) []byte
   valueType() string
}

type SliceVarint []uint64

type SliceFixed64 []uint64

type SliceFixed32 []uint32

type SliceBytes []struct {
   Message Message
   Raw []byte
}

type Number = protowire.Number

func (s SliceVarint) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   }
   return buf
}

func (s SliceFixed64) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
      buf = protowire.AppendFixed64(buf, val)
   }
   return buf
}

func (s SliceFixed32) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   }
   return buf
}

func (s SliceBytes) encode(buf []byte, num Number) []byte {
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Raw)
   }
   return buf
}

func (s SliceMessage) encode(buf []byte, num Number) []byte {
   for _, mes := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, Marshal(mes))
   }
   return buf
}

func (m Message) MarshalBinary() ([]byte, error) {
   var (
      nums []Number
      bufs []byte
   )
   for num := range mes {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      bufs = mes[num].encode(bufs, num)
   }
   return bufs, nil
}

type SliceMessage []Message

func (SliceBytes) valueType() string { return "SliceBytes" }

func (SliceFixed32) valueType() string { return "SliceFixed32" }

func (SliceFixed64) valueType() string { return "SliceFixed64" }

func (SliceVarint) valueType() string { return "SliceVarint" }

func (SliceMessage) valueType() string { return "SliceMessage" }

type typeError struct {
   Number
   in Encoder
   out Encoder
}

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

type Message map[Number]Encoder

func (m Message) varint(num Number, val uint64) error {
   if in := m[num]; in == nil {
      m[num] = SliceVarint{val}
   } else if out, ok := in.(SliceVarint); ok {
      m[num] = append(out, val)
   } else {
      return typeError{num, in, out}
   }
   return nil
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
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(buf)
         err := varint(mes, num, val)
         if err != nil {
            return nil, err
         }
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(buf)
         mes[num] = append(mes[num], Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(buf)
         mes[num] = append(mes[num], Fixed64(val))
      case protowire.BytesType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeBytes(buf)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         mes[num] = append(mes[num], val)
      }
      if err := protowire.ParseError(vLen); err != nil {
         return nil, err
      }
      buf = buf[vLen:]
   }
   return mes, nil
}
