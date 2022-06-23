package protobuf

import (
   "bufio"
   "bytes"
   "encoding/base64"
   "encoding/binary"
   "errors"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
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
      case protowire.EndGroupType: // 4
         // break would only escape switch
         return mes, nil
      case protowire.StartGroupType: // 3
         val, err = Decode(buf)
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

func (r Raw) String() string {
   if format.IsString(r) {
      return string(r)
   }
   return base64.StdEncoding.EncodeToString(r)
}

func (r Raw) MarshalText() ([]byte, error) {
   text := r.String()
   return []byte(text), nil
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

func (b Bytes) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   return protowire.AppendBytes(tag, b.Raw), nil
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

type Fixed32 uint32

func (f Fixed32) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(tag, uint32(f)), nil
}

type Fixed64 uint64

func (f Fixed64) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(tag, uint64(f)), nil
}

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
      return "", typeError{num, in, out}
   }
   return string(out.Raw), nil
}

func (m Message) GetVarint(num Number) (uint64, error) {
   in := m[num]
   out, ok := in.(Varint)
   if !ok {
      return 0, typeError{num, in, out}
   }
   return uint64(out), nil
}

func (m Message) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.BytesType)
   val, err := m.MarshalBinary()
   if err != nil {
      return nil, err
   }
   return protowire.AppendBytes(tag, val), nil
}

// we need this, so we can avoid importing
// google.golang.org/protobuf/encoding/protowire
// in other modules
type Number = protowire.Number

type Raw []byte

type Varint uint64

func (v Varint) encode(num Number) ([]byte, error) {
   tag := protowire.AppendTag(nil, num, protowire.VarintType)
   return protowire.AppendVarint(tag, uint64(v)), nil
}

func (Bytes) valueType() string { return "Bytes" }

type Encoder interface {
   encode(Number) ([]byte, error)
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
