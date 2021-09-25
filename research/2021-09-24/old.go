package html
import "github.com/tdewolff/parse/v2/html"

type Lexer struct {
   *html.Lexer
   html.TokenType
   data []byte
}

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
