package html

import (
   "fmt"
   "os"
   "testing"
)

func TestHTML(t *testing.T) {
   f, err := os.Open("index.html")
   if err != nil {
      t.Fatal(err)
   }
   defer f.Close()
   l := NewLexer(f)
   for l.NextTag("script") {
      fmt.Printf("%q\n\n", l.Bytes())
   }
}
