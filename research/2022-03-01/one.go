package hls

import (
   "io"
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

type master struct {
   bandwidth string
   uri string
}

func one(src io.Reader) []master {
   var (
      buf scanner.Scanner
      mass []master
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas master
         for buf.Scan() != '\n' {
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               mas.bandwidth = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         mas.uri = buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass
}
