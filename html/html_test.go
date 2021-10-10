package html

import (
   "os"
   "testing"
)

func TestHTML(t *testing.T) {
   file, err := os.Open("index.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   lex := NewLexer(file)
   lex.NextAttr("id", "config")
   os.Stdout.Write(lex.Bytes())
}
