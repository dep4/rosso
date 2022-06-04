package hls

import (
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func (s *Scanner) Master(base *url.URL) (*Master, error) {
   var mas Master
   mas.Streams = make(map[Stream]*url.URL)
   mas.Media = make(map[Media]*url.URL)
   for {
      s.splitWords()
      if s.Scan() == scanner.EOF {
         break
      }
      var err error
      switch s.TokenText() {
      case "EXT-X-STREAM-INF":
         var str Stream
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
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
         mas.Streams[str], err = base.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
      case "EXT-X-MEDIA":
         var med Media
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               mas.Media[med], err = scanURL(base, s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
      }
   }
   return &mas, nil
}

type Scanner struct {
   scanner.Scanner
}

func NewScanner(body io.Reader) *Scanner {
   var scan Scanner
   scan.Init(body)
   return &scan
}

func scanURL(base *url.URL, ref string) (*url.URL, error) {
   ref, err := strconv.Unquote(ref)
   if err != nil {
      return nil, err
   }
   return base.Parse(ref)
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
