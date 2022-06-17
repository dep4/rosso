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

func consumeTag(buf io.ByteReader) (Number, protowire.Type, error) {
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

type Bytes struct {
   Raw Raw // Do not embed to keep MarshalText scoped to this field
   Message
}

func (b Bytes) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   return protowire.AppendBytes(tag, b.Raw), nil
}

func (Bytes) valueType() string { return "Bytes" }

type Encoder interface {
   valueType() string
   encode(Number) ([]byte, error)
}

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

func (Encoders[T]) valueType() string {
   var value T
   return "[]" + value.valueType()
}

type Fixed32 uint32

func (f Fixed32) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(tag, uint32(f)), nil
}

func (Fixed32) valueType() string { return "Fixed32" }

type Fixed64 uint64

func (f Fixed64) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(tag, uint64(f)), nil
}

func (Fixed64) valueType() string { return "Fixed64" }

type Message map[Number]Encoder

func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Bytes:
      return value.Message
   case Message:
      return value
   }
   return nil
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

func (m Message) GetString(num Number) (string, error) {
   in := m[num]
   out, ok := in.(Bytes)
   if !ok {
      return "", getError{num, in, out}
   }
   return string(out.Raw), nil
}

func (m Message) GetVarint(num Number) (uint64, error) {
   in := m[num]
   out, ok := in.(Varint)
   if !ok {
      return 0, getError{num, in, out}
   }
   return uint64(out), nil
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

func (m Message) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

func (Message) valueType() string { return "Message" }

type Number = protowire.Number

type Raw []byte

type Varint uint64

func (v Varint) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.VarintType)
   return protowire.AppendVarint(tag, uint64(v)), nil
}

func (Varint) valueType() string { return "Varint" }

type getError struct {
   Number
   in Encoder
   out Encoder
}

func (g getError) Error() string {
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(g.Number), 10)
   b = append(b, " is "...)
   b = append(b, g.in.valueType()...)
   b = append(b, ", not "...)
   b = append(b, g.out.valueType()...)
   return string(b)
}
