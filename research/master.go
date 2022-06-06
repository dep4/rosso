package hls

import (
   "bufio"
   "bytes"
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func isMedia(s []byte) bool {
   prefix := []byte("#EXT-X-MEDIA:")
   return bytes.HasPrefix(s, prefix)
}

func isStream(s []byte) bool {
   prefix := []byte("#EXT-X-STREAM-INF:")
   return bytes.HasPrefix(s, prefix)
}

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

type Scanner struct {
   bufio *bufio.Scanner
   scanner.Scanner
}

func NewScanner(body io.Reader) Scanner {
   var scan Scanner
   scan.bufio = bufio.NewScanner(body)
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

type Streams []Stream

func (s Scanner) Master() (*Master, error) {
   var (
      err error
      mas Master
   )
   for s.bufio.Scan() {
      slice := s.bufio.Bytes()
      s.Init(bytes.NewReader(slice))
      switch {
      case isMedia(slice):
         var med Medium
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.GroupID, err = strconv.Unquote(s.TokenText())
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
            case "URI":
               s.Scan()
               s.Scan()
               med.uri, err = strconv.Unquote(s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case isStream(slice):
         var str Stream
         for s.Scan() != scanner.EOF {
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
            case "VIDEO-RANGE":
               s.Scan()
               s.Scan()
               str.VideoRange = s.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         s.bufio.Scan()
         str.uri = s.bufio.Text()
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle missing resolution
   VideoRange string // handle duplicate bandwidth
   uri string
}

func (s Stream) String() string {
   var buf []byte
   if s.Resolution != "" {
      buf = append(buf, "Resolution:"...)
      buf = append(buf, s.Resolution...)
   }
   buf = append(buf, "Bandwidth:"...)
   buf = strconv.AppendInt(buf, s.Bandwidth, 10)
   if s.Codecs != "" {
      buf = append(buf, " Codecs:"...)
      buf = append(buf, s.Codecs...)
   }
   buf = append(buf, " Range:"...)
   buf = append(buf, s.VideoRange...)
   return string(buf)
}

type Medium struct {
   Type string
   Name string
   GroupID string
   uri string
}

func (m Medium) String() string {
   var buf strings.Builder
   buf.WriteString("Type:")
   buf.WriteString(m.Type)
   buf.WriteString(" Name:")
   buf.WriteString(m.Name)
   buf.WriteString(" ID:")
   buf.WriteString(m.GroupID)
   return buf.String()
}

func (s Stream) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(s.uri)
}

func (m Medium) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(m.uri)
}
