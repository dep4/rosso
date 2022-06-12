package protobuf

import (
   "bytes"
   "encoding/base64"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

func (r Raw) MarshalText() ([]byte, error) {
   if format.IsString(r) {
      return r, nil
   }
   buf := new(bytes.Buffer)
   base64.NewEncoder(base64.StdEncoding, buf).Write(r)
   return buf.Bytes(), nil
}

type Token interface {
   appendField([]byte, Number) []byte
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
      case protowire.BytesType:
         var val []byte
         val, vLen = protowire.ConsumeBytes(data)
         embed := make(Message)
         // if this fails, then its just bytes
         embed.UnmarshalBinary(val)
         add(m, num, Bytes{val, embed})
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(data)
         add(m, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(data)
         add(m, num, Fixed64(val))
      case protowire.StartGroupType:
         var val []byte
         val, vLen = protowire.ConsumeGroup(num, data)
         embed := make(Message)
         err := embed.UnmarshalBinary(val)
         if err != nil {
            return err
         }
         add(m, num, embed)
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(data)
         add(m, num, Varint(val))
      }
      if err := protowire.ParseError(vLen); err != nil {
         return err
      }
      data = data[vLen:]
   }
   return nil
}

func (m Message) Marshal() []byte {
   var (
      data []byte
      nums []Number
   )
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      data = m[num].appendField(data, num)
   }
   return data
}

func (m Message) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

func (b Bytes) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b.Raw)
}

func (f Fixed32) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

func (f Fixed64) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

func (t Tokens[T]) appendField(buf []byte, num Number) []byte {
   for _, token := range t {
      buf = token.appendField(buf, num)
   }
   return buf
}

func (v Varint) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}
