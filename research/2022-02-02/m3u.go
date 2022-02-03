package m3u

import (
   "io"
   "strconv"
   "text/scanner"
   "unicode"
)

type Format struct {
   resolution, codecs, uri string
}

func Decode(src io.Reader, dir string) ([]Format, error) {
   var forms []Format
   var text scanner.Scanner
   text.Init(src)
   text.Whitespace = 1 << ' '
   for {
      text.IsIdentRune = func(r rune, i int) bool {
         return r == '-' || unicode.IsDigit(r) || unicode.IsLetter(r)
      }
      if text.Scan() == scanner.EOF {
         break
      }
      if text.TokenText() == "EXT-X-STREAM-INF" {
         var form Format
         for text.Scan() != '\n' {
            switch text.TokenText() {
            case "RESOLUTION":
               text.Scan()
               text.Scan()
               form.resolution = text.TokenText()
            case "CODECS":
               text.Scan()
               text.Scan()
               codec, err := strconv.Unquote(text.TokenText())
               if err != nil {
                  return nil, err
               }
               form.codecs = codec
            }
         }
         text.IsIdentRune = func(r rune, i int) bool {
            return r != '\n'
         }
         text.Scan()
         form.uri = text.TokenText()
         forms = append(forms, form)
      }
   }
   return forms, nil
}
