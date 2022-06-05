package hls

import (
   "bufio"
   "bytes"
   "io"
   "net/url"
   "strconv"
)

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
