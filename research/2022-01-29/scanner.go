package scanner

import (
   "io"
   "text/scanner"
   "unicode"
)

type format struct {
   width, height, codecs string
}

/*
Ident -2
Int -3
Float -4
Char -5
String -6
*/
func newFormat(src io.Reader) format {
   var (
      form format
      text scanner.Scanner
   )
   text.Init(src)
   text.IsIdentRune = func(r rune, i int) bool {
      return r == '#' || r == '-' || unicode.IsUpper(r)
   }
   for text.Scan() != scanner.EOF {
      switch text.TokenText() {
      case "CODECS":
         text.Scan()
         text.Scan()
         form.codecs = text.TokenText()
      case "RESOLUTION":
         text.Scan()
         text.Scan()
         form.width = text.TokenText()
         text.Scan()
         text.Scan()
         form.height = text.TokenText()
      }
   }
   return form
}
