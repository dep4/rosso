package protobuf

import (
   "fmt"
   "github.com/89z/format"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
   "sort"
)

type Bytes struct {
   Raw []byte
   Message
}

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number]Token

func Decode(src io.Reader) (Message, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   return Unmarshal(buf)
}

// Using this function strings will always be String, and []byte will always be
// Bytes. Message will be Bytes or String.
func Unmarshal(buf []byte) (Message, error) {
   if len(buf) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   mes := make(Message)
   for len(buf) >= 1 {
      num, typ, tLen := protowire.ConsumeTag(buf)
      err := protowire.ParseError(tLen)
      if err != nil {
         return nil, err
      }
      buf = buf[tLen:]
      var vLen int
      switch typ {
      case protowire.Fixed32Type:
         var val uint32
         val, vLen = protowire.ConsumeFixed32(buf)
         add(mes, num, Fixed32(val))
      case protowire.Fixed64Type:
         var val uint64
         val, vLen = protowire.ConsumeFixed64(buf)
         add(mes, num, Fixed64(val))
      case protowire.VarintType:
         var val uint64
         val, vLen = protowire.ConsumeVarint(buf)
         add(mes, num, Varint(val))
      case protowire.StartGroupType:
         var val []byte
         val, vLen = protowire.ConsumeGroup(num, buf)
         embed, err := Unmarshal(val)
         if err != nil {
            return nil, err
         }
         add(mes, num, embed)
      case protowire.BytesType:
         var val []byte
         val, vLen = protowire.ConsumeBytes(buf)
         embed, err := Unmarshal(val)
         if err != nil {
            if format.IsString(val) {
               add(mes, num, String{string(val), nil})
            } else {
               add(mes, num, Bytes{val, nil})
            }
         } else if format.IsString(val) {
            add(mes, num, String{string(val), embed})
         } else {
            add(mes, num, Bytes{val, embed})
         }
      }
      if err := protowire.ParseError(vLen); err != nil {
         return nil, err
      }
      buf = buf[vLen:]
   }
   return mes, nil
}

func (m Message) Add(num Number, val Message) {
   add(m, num, val)
}

func (m Message) AddString(num Number, val String) {
   add(m, num, val)
}

func (m Message) Get(num Number) Message {
   switch value := m[num].(type) {
   case Bytes:
      return value.Message
   case String:
      return value.Message
   }
   return nil
}

func (m Message) GetBytes(num Number) (Bytes, error) {
   return get[Bytes](m, num)
}

func (m Message) GetFixed64(num Number) (Fixed64, error) {
   return get[Fixed64](m, num)
}

func (m Message) GetMessages(num Number) []Message {
   var mes []Message
   switch value := m[num].(type) {
   case String:
      return []Message{value.Message}
   case Bytes:
      return []Message{value.Message}
   case Tokens[String]:
      for _, val := range value {
         mes = append(mes, val.Message)
      }
   case Tokens[Bytes]:
      for _, val := range value {
         mes = append(mes, val.Message)
      }
   }
   return mes
}

func (m Message) GetString(num Number) (String, error) {
   return get[String](m, num)
}

func (m Message) GetVarint(num Number) (Varint, error) {
   return get[Varint](m, num)
}

func (m Message) Marshal() []byte {
   var (
      buf []byte
      nums []Number
   )
   for num := range m {
      nums = append(nums, num)
   }
   sort.Slice(nums, func(a, b int) bool {
      return nums[a] < nums[b]
   })
   for _, num := range nums {
      buf = m[num].appendField(buf, num)
   }
   return buf
}

type Number = protowire.Number

type String struct {
   Raw string
   Message
}

type Token interface {
   appendField([]byte, Number) []byte
}

type Tokens[T Token] []T

type Varint uint64

func add[T Token](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Tokens[T]{value, val}
   case Tokens[T]:
      mes[num] = append(value, val)
   }
}

func (f Fixed32) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(in, uint32(f))
}

func (f Fixed64) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(in, uint64(f))
}

func (t Tokens[T]) appendField(in []byte, num Number) []byte {
   for _, tok := range t {
      in = tok.appendField(in, num)
   }
   return in
}

func (v Varint) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.VarintType)
   return protowire.AppendVarint(in, uint64(v))
}

func (s String) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendString(in, s.Raw)
}

func (b Bytes) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendBytes(in, b.Raw)
}

func (m Message) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendBytes(in, m.Marshal())
}

func get[T Token](mes Message, num Number) (T, error) {
   var err error
   a := mes[num]
   b, ok := a.(T)
   if !ok {
      err = fmt.Errorf("cannot unmarshal %T into field %v of type %T", a, num, b)
   }
   return b, err
}
