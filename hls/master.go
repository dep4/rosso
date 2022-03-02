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
      return r == '-' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
   }
   buf.Whitespace = 1 << ' '
}

type Master struct {
   Media []Media
   Stream []Stream
}

func NewMaster(res *http.Response) (*Master, error) {
   var (
      buf scanner.Scanner
      mas Master
   )
   buf.Init(res.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-MEDIA":
         var med Media
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "AUTOSELECT":
               buf.Scan()
               buf.Scan()
               med.AutoSelect = buf.TokenText()
            case "TYPE":
               buf.Scan()
               buf.Scan()
               med.Type = buf.TokenText()
            }
         }
         mas.Media = append(mas.Media, med)
      case "EXT-X-STREAM-INF":
         var str Stream
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               str.Bandwidth = buf.TokenText()
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               str.Resolution = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         addr, err := res.Request.URL.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         str.URI = addr.String()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}

type Media struct {
   AutoSelect string
   Type string
}

type Stream struct {
   Bandwidth string
   Resolution string
   URI string
}
