package js

import (
   "fmt"
   "github.com/89z/parse/html"
   "os"
   "testing"
)

func TestHTML(t *testing.T) {
   f, err := os.Open("index.html")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   l := html.NewLexer(f)
   for l.NextTag("script") {
      b := l.Bytes()
      fmt.Printf("BEGIN\n%s\nEND\n", b)
      s := newScanner(b)
      for s.scan() {
         fmt.Printf("%s\n---\n", s.bytes())
      }
   }
}

func TestJS(t *testing.T) {
   b, err := os.ReadFile("ig.js")
   if err != nil {
      t.Fatal(err)
   }
   s := newScanner(b)
   for s.scan() {
      fmt.Printf("%s\n---\n", s.bytes())
   }
}
