package m3u

import (
   "io"
   "strconv"
   "text/scanner"
)

type Master struct {
   Codecs string
   Resolution string
   URI string
}

func Masters(src io.Reader) ([]Master, error) {
   var (
      buf scanner.Scanner
      mass []Master
   )
   buf.Init(src)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      if buf.TokenText() == "EXT-X-STREAM-INF" {
         var mas Master
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               mas.Resolution = buf.TokenText()
            case "CODECS":
               buf.Scan()
               buf.Scan()
               codec, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               mas.Codecs = codec
            }
         }
         scanLines(&buf)
         buf.Scan()
         mas.URI = buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}
