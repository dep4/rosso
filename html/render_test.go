package html

import (
   "os"
   "strings"
   "testing"
)

func TestRender(t *testing.T) {
   l := NewLexer(strings.NewReader(s))
   err := l.Render(os.Stdout, " ")
   if err != nil {
      t.Fatal(err)
   }
}
