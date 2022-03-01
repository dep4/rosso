package hls

import (
   "net/http"
   "text/scanner"
)

func four(res *http.Response) ([]segment, error) {
   var (
      buf scanner.Scanner
      segs []segment
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
         var seg segment
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
