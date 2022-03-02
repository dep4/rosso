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
      if r == '-' || unicode.IsLetter(r) {
         return true
      }
      return i >= 1 && unicode.IsDigit(r)
   }
   buf.Whitespace = 1 << ' '
}

func one(src io.Reader) master {
   var (
      buf scanner.Scanner
      mas master
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-STREAM-INF":
         var str stream
         for buf.Scan() != '\n' {
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               str.Bandwidth = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         str.URI = buf.TokenText()
         mas.stream = append(mas.stream, str)
      case "EXT-X-MEDIA":
      }
   }
   return mas
}

type media struct {
   Name string
   Type string
}

type stream struct {
   Bandwidth string
   URI string
}

type master struct {
   media []media
   stream []stream
}
