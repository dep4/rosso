package m3u

import (
   "io"
   "strconv"
   "text/scanner"
   "unicode"
)

func scanLines(r rune, i int) bool {
   return r != '\n'
}

func scanWords(r rune, i int) bool {
   return r == '-' || unicode.IsDigit(r) || unicode.IsLetter(r)
}

type Master struct {
   Codecs string
   Resolution string
   URI string
}

func Masters(src io.Reader) ([]Master, error) {
   var (
      mass []Master
      text scanner.Scanner
   )
   text.Init(src)
   text.Whitespace = 1 << ' '
   for {
      text.IsIdentRune = scanWords
      if text.Scan() == scanner.EOF {
         break
      }
      if text.TokenText() == "EXT-X-STREAM-INF" {
         var mas Master
         for text.Scan() != '\n' {
            switch text.TokenText() {
            case "RESOLUTION":
               text.Scan()
               text.Scan()
               mas.Resolution = text.TokenText()
            case "CODECS":
               text.Scan()
               text.Scan()
               codec, err := strconv.Unquote(text.TokenText())
               if err != nil {
                  return nil, err
               }
               mas.Codecs = codec
            }
         }
         text.IsIdentRune = scanLines
         text.Scan()
         mas.URI = text.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}
