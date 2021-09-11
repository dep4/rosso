package js

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/js"
   "io"
)

type Values map[string]string

func NewValues(r io.Reader) (Values, error) {
   ast, err := js.Parse(parse.NewInput(r))
   if err != nil {
      return nil, err
   }
   v := make(Values)
   for _, iStmt := range ast.BlockStmt.List {
      eStmt, ok := iStmt.(*js.ExprStmt)
      if ok {
         bExpr, ok := eStmt.Value.(*js.BinaryExpr)
         if ok {
            y, err := bExpr.Y.JSON()
            if err != nil {
               return nil, err
            }
            v[bExpr.X.JS()] = y
         }
      }
   }
   return v, nil
}
