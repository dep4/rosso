package main

import (
   "fmt"
   "vimagination.zapto.org/javascript"
   "vimagination.zapto.org/javascript/walk"
   "vimagination.zapto.org/parser"
)

type visit struct{}

func (visit) Handle(t javascript.Type) error {
   fmt.Println(t)
   return nil
}

func main() {
   t := parser.NewStringTokeniser("var x = 'lorem ipsum';")
   m, err := javascript.ParseModule(t)
   if err != nil {
      panic(err)
   }
   var v visit
   walk.Walk(m, v)
}
