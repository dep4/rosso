package protobuf

import (
   "bytes"
   "encoding/base64"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
   "strconv"
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

func (r Raw) MarshalText() ([]byte, error) {
   if format.IsString(r) {
      return r, nil
   }
   buf := new(bytes.Buffer)
   base64.NewEncoder(base64.StdEncoding, buf).Write(r)
   return buf.Bytes(), nil
}

type Bytes struct {
   Raw Raw
   Message
}

func String(s string) Bytes {
   var out Bytes
   out.Raw = []byte(s)
   return out
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
   in := m[num]
   out, ok := in.(Bytes)
   if !ok {
      return nil, typeError{num, in, out}
   }
   return out.Raw, nil
}

func (m Message) GetFixed64(num Number) (uint64, error) {
   in := m[num]
   out, ok := in.(Fixed64)
   if !ok {
      return 0, typeError{num, in, out}
   }
   return uint64(out), nil
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

type Number = protowire.Number

type Raw []byte

type Varint uint64

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
