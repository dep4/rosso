package hls

import (
   "net/http"
   "text/scanner"
)

func two(res *http.Response) (*master, error) {
   var (
      buf scanner.Scanner
      mas master
   )
   buf.Init(res.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-MEDIA":
      case "EXT-X-STREAM-INF":
         var str stream
         for buf.Scan() != '\n' {
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               str.Bandwidth = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         addr, err := res.Request.URL.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         str.URI = addr.String()
         mas.stream = append(mas.stream, str)
      }
   }
   return &mas, nil
}
