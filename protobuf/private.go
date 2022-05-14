package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
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

func get[T Token](mes Message, num Number) (T, error) {
   var err error
   a := mes[num]
   b, ok := a.(T)
   if !ok {
      err = fmt.Errorf("cannot unmarshal %T into field %v of type %T", a, num, b)
   }
   return b, err
}

func (b Bytes) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendBytes(in, b)
}

func (f Fixed32) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(in, uint32(f))
}

func (f Fixed64) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(in, uint64(f))
}

func (m Message) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendBytes(in, m.Marshal())
}

func (s String) appendField(in []byte, num Number) []byte {
   in = protowire.AppendTag(in, num, protowire.BytesType)
   return protowire.AppendString(in, string(s))
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
