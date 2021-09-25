package js

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/js"
   "io"
)

func Parse(r io.Reader) ([]string, error) {
   in := parse.NewInput(r)
   return walk(in)
}

func ParseBytes(b []byte) ([]string, error) {
   in := parse.NewInputBytes(b)
   return walk(in)
}

func walk(r *parse.Input) ([]string, error) {
   ast, err := js.Parse(r)
   if err != nil {
      return nil, err
   }
   var vis visit
   js.Walk(&vis, ast)
   return vis.nodes, nil
}

type visit struct {
   nodes []string
}

func (v *visit) Enter(n js.INode) js.IVisitor {
   node, err := n.JSON()
   if err != nil {
      return v
   }
   v.nodes = append(v.nodes, node)
   return nil
}

func (visit) Exit(js.INode) {}
