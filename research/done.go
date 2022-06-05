package hls

import (
   "bufio"
   "io"
   "net/url"
   "strconv"
   "strings"
)

type Scanner struct {
   *bufio.Scanner
}

func NewScanner(body io.Reader) Scanner {
   var scan Scanner
   scan.Scanner = bufio.NewScanner(body)
   return scan
}

func (s Scanner) isStream() bool {
   text := s.Text()
   return strings.HasPrefix(text, "#EXT-X-STREAM-INF:")
}

func (s Scanner) isURI() bool {
   text := s.Text()
   if text == "" {
      return false
   }
   if strings.HasPrefix(text, "#") {
      return false
   }
   return true
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

type Streams []Stream

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
