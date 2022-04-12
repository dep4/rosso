package hls

import (
   "net/url"
   "strconv"
   "text/scanner"
)

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

func (s *Scanner) Master(addr *url.URL) (*Master, error) {
   var mas Master
   for {
      s.splitWords()
      if s.Scan() == scanner.EOF {
         break
      }
      var err error
      switch s.TokenText() {
      case "EXT-X-MEDIA":
         var med Media
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.GroupID, err = strconv.Unquote(s.TokenText())
            case "URI":
               s.Scan()
               s.Scan()
               med.URI, err = scanURL(s.TokenText(), addr)
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case "EXT-X-STREAM-INF":
         var str Stream
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "AUDIO":
               s.Scan()
               s.Scan()
               str.Audio, err = strconv.Unquote(s.TokenText())
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.Atoi(s.TokenText())
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
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
      }
   }
   return &mas, nil
}
