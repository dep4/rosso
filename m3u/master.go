package m3u

import (
   "io"
   "strconv"
   "text/scanner"
)

func (d Decoder) Masters(src io.Reader) ([]Master, error) {
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
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               val, err := strconv.ParseInt(buf.TokenText(), 10, 64)
               if err != nil {
                  return nil, err
               }
               mas.Bandwidth = val
            case "CODECS":
               buf.Scan()
               buf.Scan()
               val, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               mas.Codecs = val
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               mas.Resolution = buf.TokenText()
            }
         }
         scanLines(&buf)
         buf.Scan()
         mas.URI = d.Dir + buf.TokenText()
         mass = append(mass, mas)
      }
   }
   return mass, nil
}

