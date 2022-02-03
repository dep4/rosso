package m3u

import (
   "io"
   "text/scanner"
   "unicode"
)

func scanLines(r rune, i int) bool {
   return r != '\n'
}

func scanWords(r rune, i int) bool {
   return r == '-' || unicode.IsDigit(r) || unicode.IsLetter(r)
}

type Segment struct {
   Key string
   URI []string
}

func NewSegment(src io.Reader) (*Segment, error) {
   var (
      seg Segment
      text scanner.Scanner
   )
   text.Init(src)
   text.Whitespace = 1 << ' '
   for {
      text.IsIdentRune = scanWords
      if text.Scan() == scanner.EOF {
         break
      }
      switch text.TokenText() {
      case "EXT-X-KEY":
         for text.Scan() != '\n' {
            if text.TokenText() == "URI" {
               text.Scan()
               text.Scan()
               seg.Key = text.TokenText()
            }
         }
      case "EXTINF":
         text.IsIdentRune = scanLines
         text.Scan()
         seg.URI = append(seg.URI, text.TokenText())
      }
   }
   return &seg, nil
}
