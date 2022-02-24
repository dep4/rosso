package m3u

import (
   "io"
   "strconv"
   "text/scanner"
   "unicode"
)

func (d Decoder) Segments(src io.Reader) []string {
   var (
      buf scanner.Scanner
      segs []string
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXTINF" {
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         segs = append(segs, d.Dir + buf.TokenText())
      }
   }
   return segs
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

type Master struct {
   Bandwidth int64
   URI string
}

type Decoder struct {
   Dir string
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
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               val, err := strconv.ParseInt(buf.TokenText(), 10, 64)
               if err != nil {
                  return nil, err
               }
               mas.Bandwidth = val
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

