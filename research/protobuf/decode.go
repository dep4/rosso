package protobuf

import (
   "bufio"
   "bytes"
   "encoding/binary"
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "strconv"
)

func consumeBytes(buf *bufio.Reader) (Bytes, error) {
   var (
      limit io.LimitedReader
      val Bytes
   )
   n, err := binary.ReadUvarint(buf)
   if err != nil {
      return val, err
   }
   limit.N = int64(n)
   limit.R = buf
   val.Raw, err = io.ReadAll(&limit)
   if err != nil {
      return val, err
   }
   val.Message, _ = Decode(bufio.NewReader(bytes.NewReader(val.Raw)))
   return val, nil
}

func consumeFixed32(buf io.Reader) (Fixed32, error) {
   var val Fixed32
   err := binary.Read(buf, binary.LittleEndian, &val)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeFixed64(buf io.Reader) (Fixed64, error) {
   var val Fixed64
   err := binary.Read(buf, binary.LittleEndian, &val)
   if err != nil {
      return 0, err
   }
   return val, nil
}

func consumeVarint(buf io.ByteReader) (Varint, error) {
   val, err := binary.ReadUvarint(buf)
   if err != nil {
      return 0, err
   }
   return Varint(val), nil
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

func Decode(buf *bufio.Reader) (Message, error) {
   mes := make(Message)
   for {
      num, typ, err := consumeTag(buf)
      if err == io.EOF {
         return mes, nil
      } else if err != nil {
         return nil, err
      }
      var val Encoder
      switch typ {
      case protowire.VarintType: // 0
         val, err = consumeVarint(buf)
      case protowire.Fixed64Type: // 1
         val, err = consumeFixed64(buf)
      case protowire.Fixed32Type: // 5
         val, err = consumeFixed32(buf)
      case protowire.BytesType: // 2
         val, err = consumeBytes(buf)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
      add(mes, num, val)
   }
}

func add[T Encoder](mes Message, num Number, val T) error {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Encoders[T]{value, val}
   case Encoders[T]:
      mes[num] = append(value, val)
   default:
      return typeError{num, value, val}
   }
   return nil
}

type Bytes struct {
   Raw Raw
   Message
}

type Encoders[T Encoder] []T

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Encoder

// we need this, so we can avoid importing
// google.golang.org/protobuf/encoding/protowire
// in other modules
type Number = protowire.Number

type Raw []byte

type Varint uint64

func (Bytes) valueType() string { return "Bytes" }

type Encoder interface {
   valueType() string
}

func (Encoders[T]) valueType() string {
   var value T
   return "[]" + value.valueType()
}

func (Fixed32) valueType() string { return "Fixed32" }

func (Fixed64) valueType() string { return "Fixed64" }

func (Message) valueType() string { return "Message" }

func (Varint) valueType() string { return "Varint" }

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
