// HTML, HTTP and URL functions.
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
package net

import (
   "bytes"
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
   text := lex.Text()
   return string(text)
}

type Node struct {
   Attr map[string]string
   Data Text
}

func Parse(src io.Reader, tag string) []Node {
   lex := html.NewLexer(parse.NewInput(src))
   var nodes []Node
   for {
      tt, _ := lex.Next()
      if tt == html.ErrorToken {
         return nodes
      }
      if tt == html.StartTagToken && text(lex) == tag {
         attrs := make(map[string]string)
         for {
            tt, _ := lex.Next()
            if tt == html.StartTagCloseToken {
               lex.Next()
               break
            }
            if tt == html.StartTagVoidToken {
               break
            }
            if tt == html.TextToken {
               break
            }
            attrs[text(lex)] = attrVal(lex)
         }
         nodes = append(nodes, Node{
            attrs, bytes.TrimSpace(lex.Text()),
         })
      }
   }
}

type Text []byte

func (t Text) String() string {
   return string(t)
}
