package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

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

type Message map[Number]Token

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

func (m Message) GetMessages(num Number) []Message {
   var mes []Message
   switch value := m[num].(type) {
   case Bytes:
      return []Message{value.Message}
   case Tokens[Bytes]:
      for _, val := range value {
         mes = append(mes, val.Message)
      }
   }
   return mes
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

type getError struct {
   src Token
   Number
   dst Token
}

func (g getError) Error() string {
   b := new(strings.Builder)
   fmt.Fprintf(b, "cannot unmarshal %T", g.src)
   fmt.Fprintf(b, " into field %v", g.Number)
   fmt.Fprintf(b, " of type %T", g.dst)
   return b.String()
}
