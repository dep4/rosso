package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

type Raw struct {
   Bytes []byte
   String string
   Message map[Number]Encoder
}

func (Raw) get_type() string { return "Raw" }

type Bytes []byte

func (Bytes) get_type() string { return "Bytes" }

type Fixed32 uint32

func (Fixed32) get_type() string { return "Fixed32" }

type Fixed64 uint64

func (Fixed64) get_type() string { return "Fixed64" }

type Number = protowire.Number

type Encoder interface {
   get_type() string
}

type String string

func (String) get_type() string { return "String" }

type Varint uint64

func (Varint) get_type() string { return "Varint" }

type type_error struct {
   Number
   lvalue Encoder
   rvalue Encoder
}

type Slice[T Encoder] []T

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

func (t type_error) Error() string {
   var buf strings.Builder
   buf.WriteString("lvalue ")
   buf.WriteString(t.lvalue.get_type())
   buf.WriteString(" rvalue ")
   buf.WriteString(t.rvalue.get_type())
   return buf.String()
}
