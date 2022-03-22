package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
)

func add[T Token](mes Message, num Number, val T) {
   switch value := mes[num].(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = tokens[T]{value, val}
   case tokens[T]:
      mes[num] = append(value, val)
   }
}

func (m Message) appendField(b []byte, n Number) []byte {
   b = protowire.AppendTag(b, n, protowire.BytesType)
   return protowire.AppendBytes(b, m.Marshal())
}

func (s String) appendField(b []byte, n Number) []byte {
   b = protowire.AppendTag(b, n, protowire.BytesType)
   return protowire.AppendString(b, string(s))
}

func (u Uint32) appendField(b []byte, n Number) []byte {
   b = protowire.AppendTag(b, n, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(u))
}

func (u Uint64) appendField(b []byte, n Number) []byte {
   b = protowire.AppendTag(b, n, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(u))
}

type tokens[T Token] []T

func (t tokens[T]) appendField(b []byte, n Number) []byte {
   for _, tok := range t {
      b = tok.appendField(b, n)
   }
   return b
}
