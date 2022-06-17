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

func (Bytes) valueType() string { return "Bytes" }

func (Fixed32) valueType() string { return "Fixed32" }

func (Fixed64) valueType() string { return "Fixed64" }

func (Message) valueType() string { return "Message" }

func (Varint) valueType() string { return "Varint" }

type Varint []uint64

func (v Varint) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range v {
      buf = protowire.AppendTag(buf, num, protowire.VarintType)
      buf = protowire.AppendVarint(buf, val)
   }
   return buf, nil
}

type Fixed64 []uint64

func (f Fixed64) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range f {
      buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
      buf = protowire.AppendFixed64(buf, val)
   }
   return buf, nil
}

type Fixed32 []uint32

func (f Fixed32) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range f {
      buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
      buf = protowire.AppendFixed32(buf, val)
   }
   return buf, nil
}

type Bytes []struct {
   Message Message
   Raw []byte
}

func (b Bytes) encode(num Number) ([]byte, error) {
   var buf []byte
   for _, val := range b {
      buf = protowire.AppendTag(buf, num, protowire.BytesType)
      buf = protowire.AppendBytes(buf, val.Raw)
   }
   return buf, nil
}

func (m Message) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

type Message map[Number]encoder

func (m Message) MarshalBinary() ([]byte, error) {
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
