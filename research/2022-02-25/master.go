package m3u

import (
   "io"
   "path"
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
   Bandwidth int64
   URI string
}

type openFile func(string) (io.ReadCloser, error)

func Masters(src string, fn openFile) ([]Master, error) {
   file, err := fn(src)
   if err != nil {
      return nil, err
   }
   defer file.Close()
   var (
      buf scanner.Scanner
      mass []Master
   )
   buf.Init(file)
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
         // FIXME if TokenText is already absolute, this fail:
         mas.URI = path.Dir(src) + buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}
