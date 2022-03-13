package hls

import (
   "fmt"
   "io"
   "net/url"
   "strconv"
   "text/scanner"
)

func (s Stream) Format(f fmt.State, r rune) {
   if s.Resolution != "" {
      fmt.Fprint(f, "Resolution:", s.Resolution, " ")
   }
   fmt.Fprint(f, "Bandwidth:", s.Bandwidth)
   fmt.Fprint(f, " Codecs:", s.Codecs)
   if s.Audio != "" {
      fmt.Fprint(f, " Audio:", s.Audio)
   }
   if r == 'u' {
      fmt.Fprint(f, " URI:", s.URI)
   }
}

func NewMaster(addr *url.URL, body io.Reader) (*Master, error) {
   var (
      buf scanner.Scanner
      err error
      mas Master
   )
   buf.Init(body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-STREAM-INF":
         var str Stream
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "RESOLUTION":
               buf.Scan()
               buf.Scan()
               str.Resolution = buf.TokenText()
            case "BANDWIDTH":
               buf.Scan()
               buf.Scan()
               str.Bandwidth, err = strconv.Atoi(buf.TokenText())
            case "CODECS":
               buf.Scan()
               buf.Scan()
               str.Codecs, err = strconv.Unquote(buf.TokenText())
            case "AUDIO":
               buf.Scan()
               buf.Scan()
               str.Audio, err = strconv.Unquote(buf.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         scanLines(&buf)
         buf.Scan()
         str.URI, err = addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         mas.Stream = append(mas.Stream, str)
      case "EXT-X-MEDIA":
         var med Media
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "GROUP-ID":
               buf.Scan()
               buf.Scan()
               med.GroupID, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
            case "URI":
               buf.Scan()
               buf.Scan()
               ref, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               med.URI, err = addr.Parse(ref)
               if err != nil {
                  return nil, err
               }
            }
         }
         mas.Media = append(mas.Media, med)
      }
   }
   return &mas, nil
}

type Bandwidth struct {
   *Master
   Target int
}

func (b Bandwidth) Less(i, j int) bool {
   return b.distance(i) < b.distance(j)
}

func (b Bandwidth) distance(i int) int {
   diff := b.Stream[i].Bandwidth - b.Target
   if diff >= 0 {
      return diff
   }
   return -diff
}

type Master struct {
   Stream []Stream
   Media []Media
}

func (m Master) GetMedia(str Stream) *Media {
   for _, med := range m.Media {
      if med.GroupID == str.Audio {
         return &med
      }
   }
   return nil
}

func (m Master) Len() int {
   return len(m.Stream)
}

func (m Master) Swap(i, j int) {
   m.Stream[i], m.Stream[j] = m.Stream[j], m.Stream[i]
}

type Media struct {
   GroupID string
   URI *url.URL
}

type Stream struct {
   Resolution string
   Bandwidth int // handle duplicate resolution
   Codecs string // handle missing resolution
   Audio string // link to Media
   URI *url.URL
}
