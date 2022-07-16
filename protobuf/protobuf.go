package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "strconv"
)

func (self type_error) Error() string {
   get_type := func(e Encoder) string {
      if e == nil {
         return "nil"
      }
      return e.get_type()
   }
   var b []byte
   b = append(b, "field "...)
   b = strconv.AppendInt(b, int64(self.Number), 10)
   b = append(b, " is "...)
   b = append(b, get_type(self.lvalue)...)
   b = append(b, ", not "...)
   b = append(b, get_type(self.rvalue)...)
   return string(b)
}

type type_error struct {
   Number
   lvalue Encoder
   rvalue Encoder
}

type Bytes []byte

func (self Bytes) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, self)
}

func (Bytes) get_type() string { return "Bytes" }

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

type Fixed32 uint32

func (self Fixed32) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(self))
}

func (Fixed32) get_type() string { return "Fixed32" }

type Fixed64 uint64

func (self Fixed64) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(self))
}

func (Fixed64) get_type() string { return "Fixed64" }

type Number = protowire.Number

type Raw struct {
   Bytes []byte
   String string
   Message map[Number]Encoder
}

func (self Raw) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, self.Bytes)
}

func (Raw) get_type() string { return "Raw" }

type Slice[T Encoder] []T

func (self Slice[T]) encode(buf []byte, num Number) []byte {
   for _, value := range self {
      buf = value.encode(buf, num)
   }
   return buf
}

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

type String string

func (self String) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendString(buf, string(self))
}

func (String) get_type() string { return "String" }

type Varint uint64

func (self Varint) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(self))
}

func (Varint) get_type() string { return "Varint" }
