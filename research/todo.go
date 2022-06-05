package hls

import (
   "io"
   "net/url"
   "strconv"
   "text/scanner"
   "unicode"
)

func (s Scanner) Master(base *url.URL) (*Master, error) {
   var mas Master
   for scan.Scan() {
      if s.isStream() {
         var str Stream
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
            case "VIDEO-RANGE":
               s.Scan()
               s.Scan()
               str.VideoRange = s.TokenText()
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         s.splitLines()
         s.Scan()
         str.URI, err = base.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
         mas.Streams = append(mas.Streams, str)
      }
      if s.isURI() {
      }
   }
   return &mas, nil
}

func (s *Scanner) splitWords() {
   s.IsIdentRune = func(r rune, i int) bool {
      if r == '-' {
         return true
      }
      if r == '.' {
         return true
      }
      if unicode.IsDigit(r) {
         return true
      }
      if unicode.IsLetter(r) {
         return true
      }
      return false
   }
   s.Whitespace = 1 << ' '
}
