package hls

import (
   "io"
   "strings"
   "text/scanner"
   "unicode"
)

func New_Scanner(body io.Reader) Scanner {
   var scan Scanner
   scan.line.Init(body)
   scan.line.IsIdentRune = func(r rune, i int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      if r == scanner.EOF {
         return false
      }
      return true
   }
   scan.IsIdentRune = func(r rune, i int) bool {
      if r == '-' {
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
   return scan
}

type Scanner struct {
   line scanner.Scanner
   scanner.Scanner
}

type Segment struct {
   Key string
   URI []string
}

func (s Scanner) Segment() Segment {
   var (
      disco bool
      seg Segment
   )
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      switch {
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         if disco && seg.Key != "" {
            return seg
         }
         disco = false
         seg.URI = append(seg.URI, line)
      case line == "#EXT-X-DISCONTINUITY":
         disco = true
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         seg.Key = line
         seg.URI = nil
      }
   }
   return seg
}
