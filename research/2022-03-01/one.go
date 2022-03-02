package hls

import (
   "net/http"
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

func one(res *http.Response) (*master, error) {
   var (
      buf scanner.Scanner
      mas master
   )
   buf.Init(res.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-MEDIA":
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
         addr, err := res.Request.URL.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         str.URI = addr.String()
         mas.stream = append(mas.stream, str)
      }
   }
   return &mas, nil
}
