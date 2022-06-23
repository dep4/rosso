package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

func add[T Encoder](mes Message, num Number, value T) error {
   switch values := mes[num].(type) {
   case nil:
      mes[num] = value
   case T:
      mes[num] = Slice[T]{values, value}
   case Slice[T]:
      mes[num] = append(values, value)
   default:
      return type_error{values, value}
   }
   return nil
}

type Bytes []byte

func (Bytes) get_type() string { return "Bytes" }

type Encoder interface {
   get_type() string
}

type Fixed32 uint32

func (Fixed32) get_type() string { return "Fixed32" }

type Fixed64 uint64

func (Fixed64) get_type() string { return "Fixed64" }

type Number = protowire.Number

type Slice[T Encoder] []T

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

type String string

func (String) get_type() string { return "String" }

type Varint uint64

func (Varint) get_type() string { return "Varint" }

type type_error struct {
   values Encoder
   value Encoder
}

func (t type_error) Error() string {
   var buf strings.Builder
   buf.WriteString("values ")
   buf.WriteString(t.values.get_type())
   buf.WriteString(" value ")
   buf.WriteString(t.value.get_type())
   return buf.String()
}
