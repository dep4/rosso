package hls

import (
   "net/http"
   "text/scanner"
)

func two(res *http.Response) ([]master, error) {
   var (
      buf scanner.Scanner
      mass []master
   )
   buf.Init(res.Body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas master
         for buf.Scan() != '\n' {
            if buf.TokenText() == "BANDWIDTH" {
               buf.Scan()
               buf.Scan()
               mas.bandwidth = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         addr, err := res.Request.URL.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         mas.uri = addr.String()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}
