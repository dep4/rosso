package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

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

func (b Bytes) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, b.Raw)
}

func (f Fixed32) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(buf, uint32(f))
}

func (f Fixed64) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(buf, uint64(f))
}

func (m Message) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.BytesType)
   return protowire.AppendBytes(buf, m.Marshal())
}

func (t Tokens[T]) appendField(buf []byte, num Number) []byte {
   for _, token := range t {
      buf = token.appendField(buf, num)
   }
   return buf
}

func (v Varint) appendField(buf []byte, num Number) []byte {
   buf = protowire.AppendTag(buf, num, protowire.VarintType)
   return protowire.AppendVarint(buf, uint64(v))
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
