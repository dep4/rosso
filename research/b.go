package hello

type number interface {
   isNumber()
}

type Int32 int32
func (Int32) isNumber(){}

type Float32 float32
func (Float32) isNumber(){}

type object map[string]number
