package hls

import (
   "io"
   "net/url"
   "strconv"
   "strings"
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

type Media struct {
   AutoSelect string
   Type string
}

func (s Stream) String() string {
   var buf strings.Builder
   if s.Resolution != "" {
      buf.WriteString("Resolution:")
      buf.WriteString(s.Resolution)
      buf.WriteByte(' ')
   }
   buf.WriteString("Bandwidth:")
   buf.WriteString(s.Bandwidth)
   buf.WriteString(" Codecs:")
   buf.WriteString(s.Codecs)
   if s.URI != "" {
      buf.WriteString(" URI:")
      buf.WriteString(s.URI)
   }
   return buf.String()
}

func NewMaster(addr *url.URL, body io.Reader) (*Master, error) {
   var (
      buf scanner.Scanner
      err error
      mas Master
   )
   buf.Init(body)
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
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               str.Resolution = buf.TokenText()
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               str.Bandwidth = buf.TokenText()
            case "CODECS":
               buf.Scan()
               buf.Scan()
               str.Codecs, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
         scanLines(&buf)
         buf.Scan()
         addr, err = addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         str.URI = addr.String()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}


type Stream struct {
   Resolution string
   Bandwidth string
   Codecs string
   URI string
}
