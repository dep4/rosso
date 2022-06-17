package protobuf

import (
   "bufio"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strconv"
)

func consumeTag(buf io.ByteReader) (protowire.Number, protowire.Type, error) {
   tag, err := binary.ReadUvarint(buf)
   if err != nil {
      return 0, 0, err
   }
   num, typ := protowire.DecodeTag(tag)
   if num < protowire.MinValidNumber {
      return 0, 0, errors.New("invalid field number")
   }
   return num, typ, nil
}

func consumeBytes(buf *bufio.Reader) ([]byte, error) {
   n, err := binary.ReadUvarint(buf)
   if err != nil {
      return nil, err
   }
   var limit io.LimitedReader
   limit.N = int64(n)
   limit.R = buf
   return io.ReadAll(&limit)
}

func (m Message) GetMessages(num Number) []Message {
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

func (g getError) Error() string {
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(g.Number), 10)
   b = append(b, " is "...)
   b = append(b, g.in.Type()...)
   b = append(b, ", not "...)
   b = append(b, g.out.Type()...)
   return string(b)
}

func (Encoders[T]) Type() string {
   var value T
   return "[]" + value.Type()
}

func (m Message) GetVarint(num Number) (uint64, error) {
   in := m[num]
   out, ok := in.(Varint)
   if !ok {
      return 0, getError{num, in, out}
   }
   return uint64(out), nil
}

type getError struct {
   Number
   in Encoder
   out Encoder
}

func (m Message) GetString(num Number) (string, error) {
   in := m[num]
   out, ok := in.(Bytes)
   if !ok {
      return "", getError{num, in, out}
   }
   return string(out.Raw), nil
}

func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Bytes:
      return value.Message
   case Message:
      return value
   }
   return nil
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

type Number = protowire.Number

type Raw []byte

type Bytes struct {
   Raw Raw // Do not embed to keep MarshalText scoped to this field
   Message
}

type Fixed32 uint32

type Fixed64 uint64

type Varint uint64

type Message map[Number]Encoder

type Encoder interface {
   Type() string
   encode(Number) ([]byte, error)
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

func (Fixed32) Type() string { return "Fixed32" }

func (Bytes) Type() string { return "Bytes" }

func (f Fixed64) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed64Type)
   val := uint64(f)
   return protowire.AppendFixed64(tag, val), nil
}

func (Fixed64) Type() string { return "Fixed64" }

func (v Varint) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.VarintType)
   val := uint64(v)
   return protowire.AppendVarint(tag, val), nil
}

func (Varint) Type() string { return "Varint" }

func (m Message) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

func (Message) Type() string { return "Message" }

type Encoders[T Encoder] []T

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
