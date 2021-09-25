package js

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/js"
   "io"
)

type Visit struct {
   Nodes []string
}

func NewVisit(r io.Reader) (*Visit, error) {
   ast, err := js.Parse(parse.NewInput(r))
   if err != nil {
      return nil, err
   }
   vis := new(Visit)
   js.Walk(vis, ast)
   return vis, nil
}

func (v *Visit) Enter(n js.INode) js.IVisitor {
   node, err := n.JSON()
   if err != nil {
      return v
   }
   v.Nodes = append(v.Nodes, node)
   return nil
}

func (Visit) Exit(js.INode) {}
