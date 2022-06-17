package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
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
   Message map[Number]Encoder
   Raw []byte
}

type SliceMessage []map[Number]Encoder

type Number = protowire.Number

func (SliceBytes) valueType() string { return "SliceBytes" }

func (SliceFixed32) valueType() string { return "SliceFixed32" }

func (SliceFixed64) valueType() string { return "SliceFixed64" }

func (SliceVarint) valueType() string { return "SliceVarint" }

func (SliceMessage) valueType() string { return "SliceMessage" }

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

func Marshal(mes map[Number]Encoder) []byte {
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
   return bufs
}
