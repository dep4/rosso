// HTML
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
// and because "tdewolff" also has a JavaScript parser, and "golang.org" does
// not.
package html

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
   "io"
)

type lexer struct {
   *html.Lexer
}

func newLexer(r io.Reader) lexer {
   return lexer{
      html.NewLexer(parse.NewInput(r)),
   }
}

func (l lexer) tagName() string {
   return string(l.Text())
}

func (l lexer) nextTag(name string) bool {
   for {
      switch tt, _ := l.Next(); tt {
      case html.ErrorToken:
         return false
      case html.StartTagToken:
         if l.tagName() == name {
            return true
         }
      }
   }
}
