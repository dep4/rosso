package hls

import (
   "encoding/hex"
   "net/url"
   "strconv"
   "strings"
   "time"
   "unicode"
)

func scanHex(s string) ([]byte, error) {
   up := strings.ToUpper(s)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func scanDuration(s string) (time.Duration, error) {
   sec, err := strconv.ParseFloat(s, 64)
   if err != nil {
      return 0, err
   }
   return time.Duration(sec * 1000) * time.Millisecond, nil
}

func scanURL(s string, addr *url.URL) (*url.URL, error) {
   ref, err := strconv.Unquote(s)
   if err != nil {
      return nil, err
   }
   return addr.Parse(ref)
}

func (s *Scanner) splitLines() {
   s.IsIdentRune = func(r rune, i int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      return true
   }
   s.Whitespace |= 1 << '\n'
   s.Whitespace |= 1 << '\r'
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