package m3u

import (
   "io"
   "strconv"
   "text/scanner"
   "time"
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

type Decoder struct {
   Dir string
}

type Information struct {
   Runtime time.Duration
   URI string
}

type Segment struct {
   Key string
   Information []Information
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

func (d Decoder) Segment(src io.Reader) (*Segment, error) {
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
         var inf Information
         inf.URI = d.Dir + buf.TokenText()
         seg.Information = append(seg.Information, inf)
      }
   }
   return &seg, nil
}

func (d Decoder) Masters(src io.Reader) ([]Master, error) {
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
         mas.URI = d.Dir + buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}
