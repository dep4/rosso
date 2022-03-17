package main

type object map[string]interface{}

type token[T any] struct {
   object
   value T
}

func newToken[T any](m object) token[T] {
   return token[T]{object: m}
}

func (t token[T]) get(key string) token[T] {
   switch val := t.object[key].(type) {
   case object:
      t.object = val
   case T:
      t.value = val
   }
   return t
}
