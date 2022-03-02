package hls

import (
   "net/http"
   "text/scanner"
)

type segmentTwo struct {
   duration string
   uri string
}

func two(res *http.Response) ([]segmentTwo, error) {
   var (
      buf scanner.Scanner
      segs []segmentTwo
   )
   buf.Init(res.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXTINF" {
         buf.Scan()
         buf.Scan()
         var seg segmentTwo
         seg.duration = buf.TokenText()
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         addr, err := res.Request.URL.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         seg.uri = addr.String()
         segs = append(segs, seg)
      }
   }
   return segs, nil
}
