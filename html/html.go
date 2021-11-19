// HTML lexer.
//
// Why use "github.com/tdewolff/parse/v2/html" instead of
// "golang.org/x/net/html"?
//
// "go.sum" with "golang.org/x/net/html" looks like this:
//  golang.org/x/net v0.0.0-20210924151903
//  golang.org/x/sys v0.0.0-20201119102817
//  golang.org/x/sys v0.0.0-20210423082822
//  golang.org/x/term v0.0.0-20201126162022
//  golang.org/x/text v0.3.6
//  golang.org/x/tools v0.0.0-20180917221912
//
// "go.sum" with "github.com/tdewolff/parse/v2/html" looks like this:
//  github.com/tdewolff/parse/v2 v2.5.21
//  github.com/tdewolff/test v1.0.6
//
// also, if you count lines of code in non test Go files, including imported
// packages, "golang.org/x/net/html" has 8,149, while
// "github.com/tdewolff/parse/v2/html" has 1,718.
package html

import (
   "bytes"
   "encoding/json"
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
   "io"
)

func attrVal(lex *html.Lexer) string {
   attr := lex.AttrVal()
   trim := bytes.Trim(attr, `"`)
   return string(trim)
}

func text(lex *html.Lexer) string {
   tex := lex.Text()
   return string(tex)
}

type StringMap map[string]string

func NewStringMap(r io.Reader) StringMap {
   metas := make(StringMap)
   lex := html.NewLexer(parse.NewInput(r))
   for {
      tt, _ := lex.Next()
      if tt == html.ErrorToken {
         return metas
      }
      if text(lex) == "meta" {
         meta := make(StringMap)
         for {
            tt, _ := lex.Next()
            if tt == html.StartTagCloseToken || tt == html.StartTagVoidToken {
               break
            }
            meta[text(lex)] = attrVal(lex)
         }
         prop, ok := meta["property"]
         if ok {
            metas[prop] = meta["content"]
         }
      }
   }
}

func (s StringMap) Struct(val interface{}) error {
   buf, err := json.Marshal(s)
   if err != nil {
      return err
   }
   return json.Unmarshal(buf, val)
}
