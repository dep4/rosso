package hls

import (
   "net/http"
   "text/scanner"
)

type information struct {
   duration string
   uri string
}

type segment struct {
   key struct {
      method string
      uri string
   }
   inf []information
}

func newSegment(res *http.Response) (*segment, error) {
   var (
      buf scanner.Scanner
      seg segment
   )
   buf.Init(res.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXTINF":
         var inf information
         buf.Scan()
         buf.Scan()
         inf.duration = buf.TokenText()
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         addr, err := res.Request.URL.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         inf.uri = addr.String()
         seg.inf = append(seg.inf, inf)
      case "EXT-X-KEY":
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "METHOD":
               buf.Scan()
               buf.Scan()
               seg.key.method = buf.TokenText()
            case "URI":
               buf.Scan()
               buf.Scan()
               seg.key.uri = buf.TokenText()
            }
         }
      }
   }
   return &seg, nil
}
