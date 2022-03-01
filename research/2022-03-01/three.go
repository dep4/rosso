package hls

import (
   "io"
   "text/scanner"
)

type segment struct {
   duration string
   uri string
}

func segments(src io.Reader) []segment {
   var (
      scan scanner.Scanner
      segs []segment
   )
   scan.Init(src)
   for {
      scanWords(&scan)
      if scan.Scan() == scanner.EOF {
         break
      }
      if scan.TokenText() == "EXTINF" {
         scanLines(&scan)
         scan.Scan()
         scan.Scan()
         var seg segment
         seg.uri = scan.TokenText()
         segs = append(segs, seg)
      }
   }
   return segs
}
