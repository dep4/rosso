package protobuf

import (
   "bufio"
   "bytes"
   "encoding/binary"
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strings"
)

type getError struct {
   src Encoder
   Number
   dst Encoder
}

func (g getError) Error() string {
   b := new(strings.Builder)
   fmt.Fprintf(b, "cannot unmarshal %T", g.src)
   fmt.Fprintf(b, " into field %v", g.Number)
   fmt.Fprintf(b, " of type %T", g.dst)
   return b.String()
}

func (m Message) GetString(num Number) (string, error) {
   src := m[num]
   dst, ok := src.(Bytes)
   if !ok {
      return "", getError{src, num, dst}
   }
   return string(dst.Raw), nil
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

func consumeBytes(buf *bufio.Reader) ([]byte, error) {
   n, err := binary.ReadUvarint(buf)
   if err != nil {
      return nil, err
   }
   limit := io.LimitedReader{R: buf}
   limit.N = int64(n)
   return io.ReadAll(&limit)
}

func readMessage(buf *bufio.Reader) (Message, error) {
   mes := make(Message)
   for {
      num, typ, err := consumeTag(buf)
      if err == io.EOF {
         return mes, nil
      } else if err != nil {
         return nil, err
      }
      switch typ {
      case protowire.VarintType:
         val, err := binary.ReadUvarint(buf)
         if err != nil {
            return nil, err
         }
         add(mes, num, Varint(val))
      case protowire.Fixed32Type:
         var val Fixed32
         err := binary.Read(buf, binary.LittleEndian, &val)
         if err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.Fixed64Type:
         var val Fixed64
         err := binary.Read(buf, binary.LittleEndian, &val)
         if err != nil {
            return nil, err
         }
         add(mes, num, val)
      case protowire.BytesType:
         var val Bytes
         val.Raw, err = consumeBytes(buf)
         if err != nil {
            return nil, err
         }
         rd := bytes.NewReader(val.Raw)
         val.Message, _ = readMessage(bufio.NewReader(rd))
         add(mes, num, val)
      }
   }
}

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
