package main

import (
   "fmt"
   "golang.org/x/net/html"
   "strings"
)

func main() {
   r := strings.NewReader("<span class='user'>John Doe</span>")
   z := html.NewTokenizer(r)
   for {
      tt := z.Next()
      if tt == html.ErrorToken {
         break
      }
      fmt.Println(z.Token())
   }
}
