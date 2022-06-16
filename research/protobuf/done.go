package protobuf

import (
   "bufio"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

func consumeBytes(r *bufio.Reader) ([]byte, error) {
   m, err := consumeVarint(r)
   if err != nil {
      return nil, err
   }
   return io.ReadAll(io.LimitReader(r, int64(m)))
}

func consumeFixed32(r io.Reader) (uint32, error) {
   var v uint32
   err := binary.Read(r, binary.LittleEndian, &v)
   if err != nil {
      return 0, err
   }
   return v, nil
}

func consumeFixed64(r io.Reader) (uint64, error) {
   var v uint64
   err := binary.Read(r, binary.LittleEndian, &v)
   if err != nil {
      return 0, err
   }
   return v, nil
}

func consumeTag(r io.ByteReader) (protowire.Number, protowire.Type, error) {
   v, err := consumeVarint(r)
   if err != nil {
      return 0, 0, err
   }
   num, typ := protowire.DecodeTag(v)
   if num < protowire.MinValidNumber {
      return 0, 0, errors.New("invalid field number")
   }
   return num, typ, nil
}

func consumeVarint(r io.ByteReader) (uint64, error) {
   return binary.ReadUvarint(r)
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
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(data)
         add(m, num, Varint(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(data)
         add(m, num, Fixed64(val))
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(data)
         add(m, num, Fixed32(val))
      case protowire.BytesType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeBytes(data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         add(m, num, val)
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

func (e Encoders[T]) encode(num Number) ([]byte, error) {
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

func (b Bytes) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   return protowire.AppendBytes(tag, b.Raw), nil
}

func (f Fixed32) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed32Type)
   val := uint32(f)
   return protowire.AppendFixed32(tag, val), nil
}

func (f Fixed64) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed64Type)
   val := uint64(f)
   return protowire.AppendFixed64(tag, val), nil
}

func (v Varint) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.VarintType)
   val := uint64(v)
   return protowire.AppendVarint(tag, val), nil
}

func (m Message) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

func add[T Encoder](mes Message, num Number, val T) {
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
      nums []Number
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

type Bytes struct {
   Raw Raw // Do not embed to keep MarshalText scoped to this field
   Message
}

type Fixed32 uint32

type Fixed64 uint64

type Number = protowire.Number

type Raw []byte

type Varint uint64

type Encoder interface {
   encode(Number) ([]byte, error)
}

type Message map[Number]Encoder

type Encoders[T Encoder] []T
