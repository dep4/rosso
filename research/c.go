package hello

type number interface {
   isNumber()
}

type Int32 int32
func (Int32) isNumber(){}

type Float32 float32
func (Float32) isNumber(){}

type object map[string]number

func add[T number](obj object, key string, val T) {
   switch value := obj[key].(type) {
   case nil:
      obj[key] = val
   case T:
      obj[key] = []T{value, val}
   case []T:
      obj[key] = append(value, val)
   }
}
