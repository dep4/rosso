package hls

import (
   "io"
   "text/scanner"
)

type segmentThree struct {
   duration string
   uri string
}

func three(src io.Reader) []segmentThree {
   var (
      buf scanner.Scanner
      segs []segmentThree
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
         var seg segmentThree
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
