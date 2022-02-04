// M3U parser
package m3u

import (
   "io"
   "strconv"
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
      return r == '-' || unicode.IsDigit(r) || unicode.IsLetter(r)
   }
   buf.Whitespace = 1 << ' '
}

type Master struct {
   Resolution string
   Bandwidth int64
   Codecs string
   URI string
}

func (m Master) String() string {
   var buf []byte
   if m.Resolution != "" {
      buf = append(buf, "Resolution:"...)
      buf = append(buf, m.Resolution...)
      buf = append(buf, ' ')
   }
   buf = append(buf, "Bandwidth:"...)
   buf = strconv.AppendInt(buf, m.Bandwidth, 10)
   buf = append(buf, " Codecs:"...)
   buf = append(buf, m.Codecs...)
   if m.URI != "" {
      buf = append(buf, " URI:"...)
      buf = append(buf, m.URI...)
   }
   return string(buf)
}

type Scanner struct {
   Error error
   Master
   scanner.Scanner
}

func (s *Scanner) Scan() bool {
   for {
      scanWords(&s.Scanner)
      if s.Scanner.Scan() == scanner.EOF {
         return false
      }
      if s.TokenText() == "EXT-X-STREAM-INF" {
         var mas Master
         for s.Scanner.Scan() != '\n' {
            switch s.TokenText() {
            case "BANDWIDTH":
               s.Scanner.Scan()
               s.Scanner.Scan()
               val, err := strconv.ParseInt(s.TokenText(), 10, 64)
               if err != nil {
                  s.Error = err
                  return false
               }
               mas.Bandwidth = val
            case "CODECS":
               s.Scanner.Scan()
               s.Scanner.Scan()
               val, err := strconv.Unquote(s.TokenText())
               if err != nil {
                  s.Error = err
                  return false
               }
               mas.Codecs = val
            case "RESOLUTION":
               s.Scanner.Scan()
               s.Scanner.Scan()
               mas.Resolution = s.TokenText()
            }
         }
         scanLines(&s.Scanner)
         s.Scanner.Scan()
         mas.URI = s.TokenText()
         s.Master = mas
         return true
      }
   }
}

type Segment struct {
   Key string
   URI []string
}

func NewSegment(src io.Reader) (*Segment, error) {
   var (
      buf scanner.Scanner
      seg Segment
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-KEY":
         for buf.Scan() != '\n' {
            if buf.TokenText() == "URI" {
               buf.Scan()
               buf.Scan()
               key, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               seg.Key = key
            }
         }
      case "EXTINF":
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         seg.URI = append(seg.URI, buf.TokenText())
      }
   }
   return &seg, nil
}
