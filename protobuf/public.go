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
      num, typ, t_len := protowire.ConsumeTag(data)
      err := protowire.ParseError(t_len)
      if err != nil {
         return err
      }
      data = data[t_len:]
      var v_len int
      switch typ {
      case protowire.BytesType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, v_len = protowire.ConsumeBytes(data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            val.Message = nil
         }
         add(m, num, val)
      case protowire.Fixed32Type:
         var val uint32
         val, v_len = protowire.ConsumeFixed32(data)
         add(m, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, v_len = protowire.ConsumeFixed64(data)
         add(m, num, Fixed64(val))
      case protowire.VarintType:
         var val uint64
         val, v_len = protowire.ConsumeVarint(data)
         add(m, num, Varint(val))
      case protowire.StartGroupType:
         var val Bytes
         val.Message = make(Message)
         val.Raw, v_len = protowire.ConsumeGroup(num, data)
         err := val.Message.UnmarshalBinary(val.Raw)
         if err != nil {
            return err
         }
         add(m, num, val.Message)
      }
      if err := protowire.ParseError(v_len); err != nil {
         return err
      }
      data = data[v_len:]
   }
   return nil
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

type Encoder interface {
   encode(Number) ([]byte, error)
   get_type() string
}

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
