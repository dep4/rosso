package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type message map[protowire.Number]token

type token struct {
   wire protowire.Type
   value []value
}

type value struct {
   integer uint64
   bytes []byte
   message message
}

func (e Encoders[T]) encode(num protowire.Number) ([]byte, error) {
   var vals []byte
   for _, encoder := range e {
      val, err := encoder.encode(num)
      if err != nil {
         return nil, err
      }
      vals = append(vals, val...)
   }
   return vals, nil
}

func (b Bytes) encode(num protowire.Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   return protowire.AppendBytes(tag, b.Raw), nil
}

func (f Fixed32) encode(num protowire.Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed32Type)
   val := uint32(f)
   return protowire.AppendFixed32(tag, val), nil
}

func (f Fixed64) encode(num protowire.Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed64Type)
   val := uint64(f)
   return protowire.AppendFixed64(tag, val), nil
}

func (v Varint) encode(num protowire.Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.VarintType)
   val := uint64(v)
   return protowire.AppendVarint(tag, val), nil
}

func (m Message) encode(num protowire.Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

func add[T Encoder](mes Message, num protowire.Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Encoders[T]{value, val}
   case Encoders[T]:
      mes[num] = append(value, val)
   }
}

func (m Message) MarshalBinary() ([]byte, error) {
   var (
      nums []protowire.Number
      vals []byte
   )
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      val, err := m[num].encode(num)
      if err != nil {
         return nil, err
      }
      vals = append(vals, val...)
   }
   return vals, nil
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
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeBytes(data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         add(m, num, val)
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(data)
         add(m, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(data)
         add(m, num, Fixed64(val))
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(data)
         add(m, num, Varint(val))
      case protowire.StartGroupType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeGroup(num, data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            return err
         }
         add(m, num, val.Message)
      }
      if err := protowire.ParseError(vLen); err != nil {
         return err
      }
      data = data[vLen:]
   }
   return nil
}

type Bytes struct {
   Raw []byte
   Message
}

type Fixed32 uint32

type Fixed64 uint64

type Varint uint64

type Message map[protowire.Number]Encoder

func (m Message) GetMessages(num protowire.Number) []Message {
   var mes []Message
   switch value := m[num].(type) {
   case Bytes:
      return []Message{value.Message}
   case Encoders[Bytes]:
      for _, val := range value {
         mes = append(mes, val.Message)
      }
   }
   return mes
}

type Encoders[T Encoder] []T

type Encoder interface {
   encode(protowire.Number) ([]byte, error)
}
