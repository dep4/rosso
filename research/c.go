package hello

type token interface {
   isToken()
}

type Int64 int
func (Int) isToken(){}

type Float64 float64
func (Float64) isToken(){}

type object map[string]token

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
