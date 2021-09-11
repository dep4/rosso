package html

import (
   "bytes"
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
   "io"
   "strings"
   stdhtml "html"
)

// html.spec.whatwg.org/multipage/syntax.html#void-elements
var VoidElement = map[string]bool{
   "br": true,
   "img": true,
   "input": true,
   "link": true,
   "meta": true,
}

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

// developer.mozilla.org/docs/Web/API/Document/getElementsByTagName
func (l *Lexer) NextTag(name string) bool {
   for {
      switch l.TokenType, _ = l.Next(); l.TokenType {
      case html.ErrorToken:
         return false
      case html.StartTagToken:
         if l.TagName() == name {
            return true
         }
      }
   }
}

// godocs.io/golang.org/x/net/html#Render
func (l Lexer) Render(w io.Writer, indent string) error {
   var ind []byte
   b := new(bytes.Buffer)
   for {
      switch t, data := l.Next(); t {
      case html.StartTagToken:
         b.Write(ind)
         b.Write(data)
         if !VoidElement[l.TagName()] {
            ind = append(ind, indent...)
         }
      case html.AttributeToken:
         b.Write(data)
      case html.StartTagCloseToken:
         b.Write(data)
         b.WriteByte('\n')
      case html.TextToken:
         if data = bytes.TrimSpace(data); data != nil {
            b.Write(ind)
            b.Write(data)
            b.WriteByte('\n')
         }
      case html.EndTagToken:
         ind = ind[len(indent):]
         b.Write(ind)
         b.Write(data)
         b.WriteByte('\n')
      case html.ErrorToken:
         return nil
      }
      if _, err := b.WriteTo(w); err != nil {
         return err
      }
   }
}

// developer.mozilla.org/docs/Web/API/Element/tagName
func (l Lexer) TagName() string {
   return string(l.Text())
}
