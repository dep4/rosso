package hls

import (
   "io"
   "strconv"
   "text/scanner"
   "time"
   "unicode"
)

type Information struct {
   Runtime time.Duration
   URI string
}

type Segment struct {
   Key string
   Information []Information
}

type Master struct {
   Resolution string
   Bandwidth int64
   Codecs string
   URI string
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
         var inf Information
         inf.URI = d.Dir + buf.TokenText()
         seg.Information = append(seg.Information, inf)
      }
   }
   return &seg, nil
}

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
