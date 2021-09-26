package main

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/js"
)

type visit struct{}

func (visit) Exit(js.INode) {}

func (v visit) Enter(n js.INode) js.IVisitor {
   n.JSON()
   return v
}

func main() {
   ast, err := js.Parse(parse.NewInputString("var x = {y: 5 > 2}"))
   if err != nil {
      panic(err)
   }
   var vis visit
   js.Walk(vis, ast)
}
