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

func consumeTag(b io.ByteReader) (Number, protowire.Type, error) {
   tag, err := binary.ReadUvarint(b)
   if err != nil {
      return 0, 0, err
   }
   num, typ := protowire.DecodeTag(tag)
   if num < protowire.MinValidNumber {
      return 0, 0, errors.New("invalid field number")
   }
   return num, typ, nil
}

func Decode(b *bufio.Reader) (Message, error) {
   mes := make(Message)
   for {
      num, typ, err := consumeTag(b)
      if err == io.EOF {
         return mes, nil
      } else if err != nil {
         return nil, err
      }
      switch typ {
      case protowire.VarintType: // 0
         err = mes.consumeVarint(num, b)
      case protowire.Fixed64Type: // 1
         err = mes.consumeFixed64(num, b)
      case protowire.Fixed32Type: // 5
         err = mes.consumeFixed32(num, b)
      case protowire.BytesType: // 2
         err = mes.consume_raw(num, b)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
      if err != nil {
         return nil, err
      }
   }
}

type Raw struct {
   Message map[Number]Encoder
   Bytes []byte
}

type Slice[T Encoder] []T

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Encoder

// we need this, so we can avoid importing
// google.golang.org/protobuf/encoding/protowire
// in other modules
type Number = protowire.Number

type Varint uint64

func (Raw) valueType() string { return "Raw" }

type Encoder interface {
   valueType() string
}

func (Slice[T]) valueType() string {
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

func (m Message) consumeFixed32(num Number, b io.Reader) error {
   var rvalue Fixed32
   err := binary.Read(b, binary.LittleEndian, &rvalue)
   if err != nil {
      return err
   }
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Fixed32:
      m[num] = Slice[Fixed32]{lvalue, rvalue}
   case Slice[Fixed32]:
      m[num] = append(lvalue, rvalue)
   default:
      return typeError{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) consumeFixed64(num Number, b io.Reader) error {
   var rvalue Fixed64
   err := binary.Read(b, binary.LittleEndian, &rvalue)
   if err != nil {
      return err
   }
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Fixed64:
      m[num] = Slice[Fixed64]{lvalue, rvalue}
   case Slice[Fixed64]:
      m[num] = append(lvalue, rvalue)
   default:
      return typeError{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) consumeVarint(num Number, b io.ByteReader) error {
   value, err := binary.ReadUvarint(b)
   if err != nil {
      return err
   }
   rvalue := Varint(value)
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Varint:
      m[num] = Slice[Varint]{lvalue, rvalue}
   case Slice[Varint]:
      m[num] = append(lvalue, rvalue)
   default:
      return typeError{num, lvalue, rvalue}
   }
   return nil
}

func (m Message) consume_raw(num Number, b *bufio.Reader) error {
   var (
      limit io.LimitedReader
      rvalue Raw
   )
   n, err := binary.ReadUvarint(b)
   if err != nil {
      return err
   }
   limit.N = int64(n)
   limit.R = b
   rvalue.Bytes, err = io.ReadAll(&limit)
   if err != nil {
      return err
   }
   rvalue.Message, _ = Decode(bufio.NewReader(bytes.NewReader(rvalue.Bytes)))
   switch lvalue := m[num].(type) {
   case nil:
      m[num] = rvalue
   case Raw:
      m[num] = Slice[Raw]{lvalue, rvalue}
   case Slice[Raw]:
      m[num] = append(lvalue, rvalue)
   default:
      return typeError{num, lvalue, rvalue}
   }
   return nil
}
