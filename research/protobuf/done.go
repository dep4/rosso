package protobuf

import (
   "strconv"
)

func add[T Encoder](mes Message, num Number, val T) error {
   in := mes[num]
   switch out := in.(type) {
   case nil:
      mes[num] = val
   case T:
      mes[num] = Slice[T]{out, val}
   case Slice[T]:
      mes[num] = append(out, val)
   default:
      return type_error{num, in, out}
   }
   return nil
}

type Slice[T Encoder] []T

func (Slice[T]) get_type() string {
   var value T
   return "[]" + value.get_type()
}

func (s Slice[T]) encode(buf []byte, num Number) []byte {
   for _, encoder := range s {
      buf = encoder.encode(buf, num)
   }
   return buf
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
