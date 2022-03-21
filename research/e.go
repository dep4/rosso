package hello

type Int64 int
func (Int64) isToken(){}

type Float64 float64
func (Float64) isToken(){}

type token interface {
   Int64 | Float64
}

func add[T token](obj object, key string, val T) {
   switch value := obj[key].(type) {
   case nil:
      obj[key] = val
   case T:
      obj[key] = []T{value, val}
   case []T:
      obj[key] = append(value, val)
   }
}

type object map[string]any
