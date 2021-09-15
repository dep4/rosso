package html

import (
   "bytes"
   "github.com/tdewolff/parse/v2/html"
   "io"
)

// html.spec.whatwg.org/multipage/syntax.html#void-elements
var VoidElement = map[string]bool{
   "br": true,
   "img": true,
   "input": true,
   "link": true,
   "meta": true,
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
