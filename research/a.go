package hello

type object map[string]any

func add[T any](obj object, key string, val T) {
   switch value := obj[key].(type) {
   case nil:
      obj[key] = val
   case T:
      obj[key] = []T{value, val}
   case []T:
      obj[key] = append(value, val)
   }
}
