package hls

import (
   "io"
   "text/scanner"
)

type segment struct {
   duration string
   uri string
}

func Five(src io.Reader) []segment {
   var (
      buf scanner.Scanner
      segs []segment
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXTINF" {
         buf.Scan()
         buf.Scan()
         var seg segment
         seg.duration = buf.TokenText()
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         seg.uri = buf.TokenText()
         segs = append(segs, seg)
      }
   }
   return segs
}
