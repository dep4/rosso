package main

import (
   "fmt"
   "github.com/tdewolff/parse/v2"
   "github.com/tdewolff/parse/v2/html"
)

func main() {
   in := parse.NewInputString("<span class='user'>John Doe</span>")
   l := html.NewLexer(in)
   for {
      tt, data := l.Next()
      if tt == html.ErrorToken {
         break
      }
      fmt.Printf("%v %s\n", tt, data)
   }
}
