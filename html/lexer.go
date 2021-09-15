package html

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
   "io"
   "strings"
   stdhtml "html"
)

// godocs.io/github.com/tdewolff/parse/v2/html#Lexer
type Lexer struct {
   *html.Lexer
   html.TokenType
   data []byte
   attr map[string]string
}

// godocs.io/github.com/tdewolff/parse/v2/html#NewLexer
func NewLexer(r io.Reader) Lexer {
   return Lexer{
      Lexer: html.NewLexer(parse.NewInput(r)),
   }
}

// developer.mozilla.org/docs/Web/API/Node/textContent
func (l *Lexer) Bytes() []byte {
   for {
      switch l.TokenType {
      case html.ErrorToken:
         return nil
      case html.TextToken:
         return l.data
      }
      l.TokenType, l.data = l.Next()
   }
}

// developer.mozilla.org/docs/Web/API/Element/getAttribute
func (l Lexer) GetAttr(key string) string {
   val := l.attr[key]
   return stdhtml.UnescapeString(strings.Trim(val, `'"`))
}

// developer.mozilla.org/docs/Web/API/Element/hasAttribute
func (l Lexer) HasAttr(key string) bool {
   _, ok := l.attr[key]
   return ok
}

// developer.mozilla.org/docs/Web/API/Document/getElementsByClassName
func (l *Lexer) NextAttr(key, val string) bool {
   for {
      switch l.TokenType, _ = l.Next(); l.TokenType {
      case html.ErrorToken:
         return false
      case html.StartTagToken:
         l.attr = make(map[string]string)
      case html.AttributeToken:
         l.attr[string(l.Text())] = string(l.AttrVal())
      case html.StartTagCloseToken, html.StartTagVoidToken:
         if l.HasAttr(key) && l.GetAttr(key) == val {
            return true
         }
      }
   }
}
