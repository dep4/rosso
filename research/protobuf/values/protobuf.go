package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Number = protowire.Number

type encoder interface {
   encode(Number) ([]byte, error)
   valueType() string
}

type SliceVarint []uint64

type SliceFixed64 []uint64

type SliceFixed32 []uint32

func (SliceBytes) valueType() string { return "SliceBytes" }

func (SliceFixed32) valueType() string { return "SliceFixed32" }

func (SliceFixed64) valueType() string { return "SliceFixed64" }

func (SliceVarint) valueType() string { return "SliceVarint" }

type SliceMessage []map[Number]encoder

func (SliceMessage) valueType() string { return "SliceMessage" }

func (s SliceVarint) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   }
   return buf, nil
}

func (s SliceFixed64) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
      buf = protowire.AppendFixed64(buf, val)
   }
   return buf, nil
}

func (s SliceFixed32) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   }
   return buf, nil
}

type SliceBytes []struct {
   Message map[Number]encoder
   Raw []byte
}

func (s SliceBytes) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range s {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Raw)
   }
   return buf, nil
}

func (m SliceMessage) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

func (m SliceMessage) MarshalBinary() ([]byte, error) {
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
      buf, err := m[num].encode(num)
      if err != nil {
         return nil, err
      }
      bufs = append(bufs, buf...)
   }
   return bufs, nil
}

func (m Message) UnmarshalBinary(data []byte) error {
   if len(data) == 0 {
      return io.ErrUnexpectedEOF
   }
   for len(data) >= 1 {
      num, typ, tLen := protowire.ConsumeTag(data)
      err := protowire.ParseError(tLen)
      if err != nil {
         return err
      }
      data = data[tLen:]
      var vLen int
      switch typ {
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(data)
         m[num] = append(m[num], Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(data)
         m[num] = append(m[num], Fixed64(val))
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(data)
         m[num] = append(m[num], Varint(val))
      case protowire.BytesType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeBytes(data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         m[num] = append(m[num], val)
      }
      if err := protowire.ParseError(vLen); err != nil {
         return err
      }
      data = data[vLen:]
   }
   return nil
}
