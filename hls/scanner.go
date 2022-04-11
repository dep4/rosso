package hls

import (
   "encoding/hex"
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func scanHex(s string) ([]byte, error) {
   s = strings.TrimPrefix(s, "0x")
   return hex.DecodeString(s)
}

func scanURL(s string, addr *url.URL) (*url.URL, error) {
   ref, err := strconv.Unquote(s)
   if err != nil {
      return nil, err
   }
   return addr.Parse(ref)
}

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
