package main
import "fmt"

type Int int
func (Int) isNumber(){}

type Ints []int
func (Ints) isNumber(){}

type number interface {
   isNumber()
}

func main() {
   obj := map[string]number{
      "one": Int(1),
      "two": Int(2),
   }
   fmt.Println(obj)
}
