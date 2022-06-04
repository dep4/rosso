package hls

import (
   "io"
   "net/url"
   "strconv"
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

func (m Media) String() string {
   var buf strings.Builder
   buf.WriteString("Type:")
   buf.WriteString(m.Type)
   buf.WriteString(" Name:")
   buf.WriteString(m.Name)
   return buf.String()
}

func (s Stream) String() string {
   var buf []byte
   if s.Resolution != "" {
      buf = append(buf, "Resolution:"...)
      buf = append(buf, s.Resolution...)
      buf = append(buf, ' ')
   }
   buf = append(buf, "Bandwidth:"...)
   buf = strconv.AppendInt(buf, s.Bandwidth, 10)
   if s.Codecs != "" {
      buf = append(buf, " Codecs:"...)
      buf = append(buf, s.Codecs...)
   }
   return string(buf)
}

type Master struct {
   Streams map[Stream]*url.URL
   Media map[Media]*url.URL
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle missing resolution
}

type Media struct {
   Name string
   Type string
}
