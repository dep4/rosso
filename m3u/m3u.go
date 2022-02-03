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

func Masters(src io.Reader, dir string) ([]Master, error) {
   var (
      buf scanner.Scanner
      mass []Master
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas Master
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               val, err := strconv.ParseInt(buf.TokenText(), 10, 64)
               if err != nil {
                  return nil, err
               }
               mas.Bandwidth = val
            case "CODECS":
               buf.Scan()
               buf.Scan()
               val, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               mas.Codecs = val
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               mas.Resolution = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         mas.URI = dir + buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}

func (m Master) String() string {
   buf := []byte("Resolution:")
   buf = append(buf, m.Resolution...)
   buf = append(buf, " Bandwidth:"...)
   buf = strconv.AppendInt(buf, m.Bandwidth, 10)
   buf = append(buf, " Codecs:"...)
   buf = append(buf, m.Codecs...)
   return string(buf)
}

type Segment struct {
   Key string
   URI []string
}

func NewSegment(src io.Reader, dir string) (*Segment, error) {
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
         seg.URI = append(seg.URI, dir + buf.TokenText())
      }
   }
   return &seg, nil
}
