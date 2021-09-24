package html

import (
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
   "io"
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
