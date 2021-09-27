package main

import (
   "fmt"
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/js"
)

type visit struct{}

func (visit) Exit(js.INode) {}

func (v visit) Enter(n js.INode) js.IVisitor {
   fmt.Println(n)
   return v
}

func main() {
   ast, err := js.Parse(parse.NewInputString("var x = 'lorem ipsum';"))
   if err != nil {
      panic(err)
   }
   var vis visit
   js.Walk(vis, ast)
}
