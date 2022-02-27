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
