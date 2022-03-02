package hls

import (
   "io"
   "net/url"
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
            case "CODECS":
               buf.Scan()
               buf.Scan()
               str.Codecs, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               str.Bandwidth, err = strconv.ParseInt(buf.TokenText(), 10, 64)
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
   Bandwidth int64
   Codecs string
   Resolution string
   URI string
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
   buf = append(buf, " Codecs:"...)
   buf = append(buf, s.Codecs...)
   if s.URI != "" {
      buf = append(buf, " URI:"...)
      buf = append(buf, s.URI...)
   }
   return string(buf)
}
