package main

import (
   "fmt"
   "github.com/robertkrimen/otto/ast"
   "github.com/robertkrimen/otto/parser"
)

type visit struct{}

func (visit) Exit(ast.Node) {}

func (v visit) Enter(n ast.Node) ast.Visitor {
   fmt.Printf("%+v\n", n)
   return v
}

func main() {
   pro, err := parser.ParseFile(nil, "", "var x = [10,11];", 0)
   if err != nil {
      panic(err)
   }
   var vis visit
   ast.Walk(vis, pro)
}
