package js

import (
   "fmt"
   "github.com/89z/parse/html"
   "os"
   "testing"
)

func TestJS(t *testing.T) {
   f, err := os.Open("index.html")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   l := html.NewLexer(f)
   for l.NextTag("script") {
      s := newScanner(l.Bytes())
      for s.scan() {
         fmt.Printf("%s\n---\n", s.bytes())
      }
   }
}
