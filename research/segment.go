package m3u

import (
   "io"
   "strconv"
   "text/scanner"
)

type Segment struct {
   Key string
   URI []string
}

func (d Decoder) Segment(src io.Reader) (*Segment, error) {
   var (
      buf scanner.Scanner
      seg Segment
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-KEY":
         for buf.Scan() != '\n' {
            if buf.TokenText() == "URI" {
               buf.Scan()
               buf.Scan()
               key, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               seg.Key = key
            }
         }
      case "EXTINF":
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         seg.URI = append(seg.URI, d.Dir + buf.TokenText())
      }
   }
   return &seg, nil
}
