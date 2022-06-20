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

func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, t_len := protowire.ConsumeTag(buf)
      err := protowire.ParseError(t_len)
      if err != nil {
         return nil, err
      }
      buf = buf[t_len:]
      // until we have exhaustive switch, we need an extra variable here
      var v_len int
      switch typ {
      case protowire.BytesType:
         var val Bytes
         val.Raw, v_len = protowire.ConsumeBytes(buf)
         val.Message, _ = Unmarshal(val.Raw)
         add(mes, num, val)
      case protowire.Fixed32Type:
         var val uint32
         val, v_len = protowire.ConsumeFixed32(buf)
         add(mes, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, v_len = protowire.ConsumeFixed64(buf)
         add(mes, num, Fixed64(val))
      case protowire.VarintType:
         var val uint64
         val, v_len = protowire.ConsumeVarint(buf)
         add(mes, num, Varint(val))
      case protowire.StartGroupType:
         var val Bytes
         val.Raw, v_len = protowire.ConsumeGroup(num, buf)
         val.Message, err = Unmarshal(val.Raw)
         if err != nil {
            return nil, err
         }
         add(mes, num, val.Message)
      }
      if err := protowire.ParseError(v_len); err != nil {
         return nil, err
      }
      buf = buf[v_len:]
   }
   return mes, nil
}

func (e Encoders[T]) encode(buf []byte, num Number) []byte {
   for _, encoder := range e {
      buf = encoder.encode(buf, num)
   }
   return buf
}

func (m Message) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

func (m Message) Marshal() []byte {
   var (
      nums []Number
      buf []byte
   )
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      buf = m[num].encode(buf, num)
   }
   return buf
}

func (v Varint) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}

func (f Fixed64) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

func (f Fixed32) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

func (b Bytes) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b.Raw)
}

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

func (r Raw) MarshalText() ([]byte, error) {
   if format.Is_String(r) {
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

func (m Message) Add_String(num Number, val string) {
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

func (m Message) Get_Bytes(num Number) ([]byte, error) {
   in := m[num]
   out, ok := in.(Bytes)
   if !ok {
      return nil, type_error{num, in, out}
   }
   return out.Raw, nil
}

func (m Message) Get_Fixed64(num Number) (uint64, error) {
   in := m[num]
   out, ok := in.(Fixed64)
   if !ok {
      return 0, type_error{num, in, out}
   }
   return uint64(out), nil
}

func (m Message) Get_String(num Number) (string, error) {
   in := m[num]
   out, ok := in.(Bytes)
   if !ok {
      return "", type_error{num, in, out}
   }
   return string(out.Raw), nil
}

func (m Message) Get_Varint(num Number) (uint64, error) {
   in := m[num]
   out, ok := in.(Varint)
   if !ok {
      return 0, type_error{num, in, out}
   }
   return uint64(out), nil
}

type Number = protowire.Number

type Raw []byte

type Varint uint64

type Message map[Number]Encoder

func (m Message) Get_Messages(num Number) []Message {
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

func (Bytes) get_type() string { return "Bytes" }

func (Encoders[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

func (Fixed32) get_type() string { return "Fixed32" }

func (Fixed64) get_type() string { return "Fixed64" }

func (Message) get_type() string { return "Message" }

func (Varint) get_type() string { return "Varint" }

type type_error struct {
   Number
   in Encoder
   out Encoder
}

func (t type_error) Error() string {
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(t.Number), 10)
   b = append(b, " is "...)
   b = append(b, t.in.get_type()...)
   b = append(b, ", not "...)
   b = append(b, t.out.get_type()...)
   return string(b)
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
