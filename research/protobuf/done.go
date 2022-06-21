package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

type Number = protowire.Number

type type_error struct {
   Number
   in Encoder
   out Encoder
}

type Encoder interface {
   encode([]byte, Number) []byte
   get_type() string
}

type Varint uint64

func (Varint) get_type() string { return "Varint" }

func (v Varint) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
}

type Fixed32 uint32

func (Fixed32) get_type() string { return "Fixed32" }

func (f Fixed32) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

type Fixed64 uint64

func (Fixed64) get_type() string { return "Fixed64" }

func (f Fixed64) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

type Bytes []byte

func (Bytes) get_type() string { return "Bytes" }

func (b Bytes) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b)
}

type Message map[Number]Encoder

func (Message) get_type() string { return "Message" }

func (m Message) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

type Raw struct {
   Bytes []byte
   Message Message
}

func (Raw) get_type() string { return "Raw" }

func (r Raw) encode(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, r.Bytes)
}
