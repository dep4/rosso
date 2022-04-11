package hls

import (
   "fmt"
   "net/url"
   "strconv"
   "text/scanner"
   "time"
)

type Bandwidth struct {
   *Master
   Target int
}

func (b Bandwidth) Less(i, j int) bool {
   distance := func(k int) int {
      diff := b.Stream[k].Bandwidth - b.Target
      if diff >= 0 {
         return diff
      }
      return -diff
   }
   return distance(i) < distance(j)
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

func (s Stream) Format(f fmt.State, verb rune) {
   if s.Resolution != "" {
      fmt.Fprint(f, "Resolution:", s.Resolution, " ")
   }
   fmt.Fprint(f, "Bandwidth:", s.Bandwidth)
   fmt.Fprint(f, " Codecs:", s.Codecs)
   if verb == 'a' {
      fmt.Fprint(f, " Audio:", s.Audio)
      fmt.Fprint(f, " URI:", s.URI)
   }
}

type Information struct {
   time.Duration
   IV []byte
   URI *url.URL
}

type Segment struct {
   Key *url.URL
   Info []Information
}

func (s *Scanner) Master(addr *url.URL) (*Master, error) {
   var (
      err error
      mas Master
   )
   for {
      s.splitWords()
      if s.Scan() == scanner.EOF {
         break
      }
      switch s.TokenText() {
      case "EXT-X-STREAM-INF":
         var str Stream
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.Atoi(s.TokenText())
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            case "AUDIO":
               s.Scan()
               s.Scan()
               str.Audio, err = strconv.Unquote(s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         s.splitLines()
         s.Scan()
         str.URI, err = addr.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
         mas.Stream = append(mas.Stream, str)
      case "EXT-X-MEDIA":
         var med Media
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.GroupID, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            case "URI":
               s.Scan()
               s.Scan()
               ref, err := strconv.Unquote(s.TokenText())
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
