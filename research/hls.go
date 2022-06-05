package hls

import (
   "bufio"
   "bytes"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "text/scanner"
)

func (s Stream) Format(f fmt.State, verb rune) {
   if s.Resolution != "" {
      fmt.Fprint(f, "Resolution:", s.Resolution, " ")
   }
   fmt.Fprint(f, "Bandwidth:", s.Bandwidth)
   if s.Codecs != "" {
      fmt.Fprint(f, " Codecs:", s.Codecs)
   }
   if verb == 'a' {
      fmt.Fprint(f, " Range:", s.VideoRange)
      fmt.Fprint(f, " URI:", s.URI)
   }
}

func (s Scanner) Stream(base *url.URL) (*Stream, error) {
   var t scanner.Scanner
   t.Init(bytes.NewReader(s.Bytes()))
   /*
   t.Whitespace = 1 << ' '
   t.IsIdentRune = func(r rune, i int) bool {
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
   */
   var (
      err error
      str Stream
   )
   for t.Scan() != scanner.EOF {
      switch t.TokenText() {
      case "BANDWIDTH":
         t.Scan()
         t.Scan()
         str.Bandwidth, err = strconv.ParseInt(t.TokenText(), 10, 64)
      case "CODECS":
         t.Scan()
         t.Scan()
         str.Codecs, err = strconv.Unquote(t.TokenText())
      case "RESOLUTION":
         t.Scan()
         t.Scan()
         str.Resolution = t.TokenText()
      case "VIDEO-RANGE":
         t.Scan()
         t.Scan()
         str.VideoRange = t.TokenText()
      }
      if err != nil {
         return nil, err
      }
   }
   s.Scan()
   str.URI, err = base.Parse(s.Text())
   if err != nil {
      return nil, err
   }
   return &str, nil
}
type Scanner struct {
   *bufio.Scanner
}

func NewScanner(body io.Reader) Scanner {
   var scan Scanner
   scan.Scanner = bufio.NewScanner(body)
   return scan
}

func scanURL(base *url.URL, raw string) (*url.URL, error) {
   ref, err := strconv.Unquote(raw)
   if err != nil {
      return nil, err
   }
   return base.Parse(ref)
}

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

type Streams []*Stream

type Medium struct {
   Type string
   Name string
   GroupID string
   URI *url.URL
}

type Stream struct {
   Resolution string
   VideoRange string // handle duplicate bandwidth
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle missing resolution
   URI *url.URL
}

func (s Scanner) Master(base *url.URL) (*Master, error) {
   var (
      mas Master
      prefix = []byte("#EXT-X-STREAM-INF:")
   )
   for s.Scan() {
      slice := s.Bytes()
      if bytes.HasPrefix(slice, prefix) {
         stream, err := s.Stream(base)
         if err != nil {
            return nil, err
         }
         mas.Streams = append(mas.Streams, stream)
      }
   }
   return &mas, nil
}
