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
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
   "io"
)

type Lexer struct {
   *html.Lexer
}

func NewLexer(r io.Reader) Lexer {
   inp := parse.NewInput(r)
   return Lexer{
      html.NewLexer(inp),
   }
}

func (l Lexer) AttrVal() []byte {
   val := l.Lexer.AttrVal()
   return bytes.Trim(val, `'"`)
}

// Keep going until we reach "Text", "EndTag" (</script>), "StartTagVoid" (/>)
// or "StartTag" (<script>). Typically this method would not be used with void
// elements, as they have no children. However if used with a void element, and
// a text node immediately follows, it will be returned. Ideally "nil" would be
// returned, but that would require maintaining a list of all void elements.
func (l Lexer) Bytes() []byte {
   for {
      switch tt, data := l.Next(); tt {
      case html.ErrorToken, html.EndTagToken:
         return nil
      case html.TextToken, html.StartTagVoidToken, html.StartTagToken:
         return data
      }
   }
}

func (l Lexer) NextAttr(key, val string) bool {
   for {
      switch tt, _ := l.Next(); tt {
      case html.ErrorToken:
         return false
      case html.AttributeToken:
         if string(l.Text()) == key && string(l.AttrVal()) == val {
            return true
         }
      }
   }
}
