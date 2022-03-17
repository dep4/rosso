package protobuf

type token[T Message | []byte | string | uint64] struct {
   value T
}

func (t token[T]) add(m Message, num Number, typ Type) {
   key := Tag{num, typ}
   switch value := m[key].(type) {
   case nil:
      m[key] = t.value
   case T:
      m[key] = []T{value, t.value}
   case []T:
      m[key] = append(value, t.value)
   }
}
