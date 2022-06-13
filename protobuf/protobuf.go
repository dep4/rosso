package protobuf

import (
   "bytes"
   "encoding/base64"
   "fmt"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strings"
)

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

type Encoder interface {
   encode(Number) ([]byte, error)
}

type Message map[Number]Encoder

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

type Encoders[T Encoder] []T

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

type getError struct {
   src Encoder
   Number
   dst Encoder
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
      case protowire.StartGroupType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, vLen = protowire.ConsumeGroup(num, data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            return err
         }
         add(m, num, val.Message)
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

func (r Raw) MarshalText() ([]byte, error) {
   if format.IsString(r) {
      return r, nil
   }
   buf := new(bytes.Buffer)
   base64.NewEncoder(base64.StdEncoding, buf).Write(r)
   return buf.Bytes(), nil
}

type Bytes struct {
   Raw Raw // Do not embed to keep MarshalText scoped to this field
   Message
}

func String(s string) Bytes {
   var dst Bytes
   dst.Raw = []byte(s)
   return dst
}

type Fixed32 uint32

type Fixed64 uint64

func (m Message) Add(num Number, val Message) {
   add(m, num, val)
}

func (m Message) AddString(num Number, val string) {
   add(m, num, String(val))
}

// Check Bytes for Unmarshaled Messages, check Message for manually constructed
// Messages.
func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Bytes:
      return value.Message
   case Message:
      return value
   }
   return nil
}

func (m Message) GetBytes(num Number) ([]byte, error) {
   src := m[num]
   dst, ok := src.(Bytes)
   if !ok {
      return nil, getError{src, num, dst}
   }
   return dst.Raw, nil
}

func (m Message) GetFixed64(num Number) (uint64, error) {
   src := m[num]
   dst, ok := src.(Fixed64)
   if !ok {
      return 0, getError{src, num, dst}
   }
   return uint64(dst), nil
}

func (m Message) GetString(num Number) (string, error) {
   src := m[num]
   dst, ok := src.(Bytes)
   if !ok {
      return "", getError{src, num, dst}
   }
   return string(dst.Raw), nil
}

func (m Message) GetVarint(num Number) (uint64, error) {
   src := m[num]
   dst, ok := src.(Varint)
   if !ok {
      return 0, getError{src, num, dst}
   }
   return uint64(dst), nil
}

type Number = protowire.Number

type Raw []byte

type Varint uint64

func (g getError) Error() string {
   b := new(strings.Builder)
   fmt.Fprintf(b, "cannot unmarshal %T", g.src)
   fmt.Fprintf(b, " into field %v", g.Number)
   fmt.Fprintf(b, " of type %T", g.dst)
   return b.String()
}
