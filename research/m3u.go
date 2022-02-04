package m3u

import (
   "io"
   "strconv"
   "text/scanner"
   "unicode"
)

func scanLines(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r != '\n'
   }
   buf.Whitespace = 1 << '\n'
}

func scanWords(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r == '-' || unicode.IsDigit(r) || unicode.IsLetter(r)
   }
   buf.Whitespace = 1 << ' '
}

type Master struct {
   Resolution string
   Bandwidth int64
   Codecs string
   URI string
}

type Scanner struct {
   Error error
   Master
   io.Reader
}

func (s *Scanner) Scan() bool {
   var buf scanner.Scanner
   buf.Init(s)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         return false
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas Master
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               val, err := strconv.ParseInt(buf.TokenText(), 10, 64)
               if err != nil {
                  s.Error = err
                  return false
               }
               mas.Bandwidth = val
            case "CODECS":
               buf.Scan()
               buf.Scan()
               val, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  s.Error = err
                  return false
               }
               mas.Codecs = val
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               mas.Resolution = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         mas.URI = buf.TokenText()
         s.Master = mas
         return true
      }
   }
}
