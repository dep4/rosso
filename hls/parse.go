package hls

import (
   "encoding/hex"
   "io"
   "strings"
   "text/scanner"
   "unicode"
)

type Scanner struct {
   scanner.Scanner
}

func NewScanner(body io.Reader) *Scanner {
   var scan Scanner
   scan.Init(body)
   return &scan
}

func (s *Scanner) splitLines() {
   s.Whitespace |= 1 << '\n'
   s.Whitespace |= 1 << '\r'
   s.IsIdentRune = func(r rune, i int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      return true
   }
}

func (s *Scanner) splitWords() {
   s.Whitespace = 1 << ' '
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
}

////////////////////////////////////////////////////////////////////////////////

func hexDecode(in string) ([]byte, error) {
   in = strings.TrimPrefix(in, "0x")
   return hex.DecodeString(in)
}
