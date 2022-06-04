package hls

import (
   "net/url"
   "strconv"
   "text/scanner"
)

func (s *Scanner) Master(base *url.URL) (*Master, error) {
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
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               med.URI, err = scanURL(s.TokenText(), base)
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
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
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
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
         str.URI, err = base.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}

// Provide a name such as "English"
func (m Master) Audio(name string) *Media {
   for _, med := range m.Media {
      if med.Type == "AUDIO" && med.Name == name {
         return &med
      }
   }
   return nil
}

func (m Master) Stream(bandwidth int64) *Stream {
   distance := func(s *Stream) int64 {
      if s.Bandwidth > bandwidth {
         return s.Bandwidth - bandwidth
      }
      return bandwidth - s.Bandwidth
   }
   var dst *Stream
   for i, src := range m.Streams {
      if dst == nil || distance(&src) < distance(dst) {
         dst = &m.Streams[i]
      }
   }
   return dst
}
