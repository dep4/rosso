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

func (l lexer) nextTag(name string) bool {
   for {
      // the second return value look like "<script"
      switch tt, _ := l.Next(); tt {
      case html.ErrorToken:
         return false
      case html.StartTagToken:
         if string(l.Text()) == name {
            return true
         }
      }
   }
}

// Keep going until we reach "Text", "EndTag", "StartTagVoid" or "StartTag". If
// the current element is void, such as <meta>, this might produce unexpected
// result. This is a compromise, as the fix would be to maintain a list of all
// void elements.
func (l lexer) nextText() bool {
   for {
      switch tt, _ := l.Next(); tt {
      case html.ErrorToken:
         return false
      case
      html.TextToken,
      html.EndTagToken,
      html.StartTagVoidToken,
      html.StartTagToken:
         return true
      }
   }
}
