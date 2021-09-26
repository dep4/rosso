package main

import (
   "github.com/robertkrimen/otto/parser"
)

const src = `
window.__additionalDataLoaded(
   true,
   9.9,
   'extra',
   [9,8],
   {"shortcode_media":9},
   null
);
`

func main() {
   _, err := parser.ParseFile(nil, "", src, 0)
   if err != nil {
      panic(err)
   }
}
