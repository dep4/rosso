package hls

import (
   "io"
   "net/http"
   "strconv"
   "text/scanner"
   "unicode"
)

func one(req *http.Request) (*response, error) {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   return &response{res}, nil
}

func two(rc io.ReadCloser) *response {
   return nil
}

type response struct {
   *http.Response
}

// r.Request.URL
func (r response) masters() ([]master, error) {
   defer r.Body.Close()
   var (
      buf scanner.Scanner
      mass []master
   )
   buf.Init(r.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas master
         for buf.Scan() != '\n' {
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               val, err := strconv.ParseInt(buf.TokenText(), 10, 64)
               if err != nil {
                  return nil, err
               }
               mas.bandwidth = val
            }
         }
         scanLines(&buf)
         buf.Scan()
         // FIXME add ResolveReference
         mas.uri = buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}

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

type master struct {
   bandwidth int64
   uri string
}
