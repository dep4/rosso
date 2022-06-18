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
      buf = protowire.AppendBytes(buf, mes.MarshalBinary())
   }
   return buf
}

// FIXME should return error
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

type Message map[Number]Encoder

type SliceMessage []Message

type Bytes struct {
   Raw []byte
   Message Message
}

type SliceBytes []Bytes
