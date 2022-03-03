package hls

import (
   "io"
   "net/url"
   "strconv"
   "text/scanner"
)

func (m Master) URIs(f func(Stream) bool) []string {
   for _, str := range m.Stream {
      if f(str) {
         uris := []string{str.URI}
         for _, med := range m.Media {
            if med.GroupID == str.Audio {
               uris = append(uris, med.URI)
            }
         }
         return uris
      }
   }
   return nil
}

type Master struct {
   Stream []Stream
   Media []Media
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle audio only
   Audio string // link to Media
   URI string
}

type Media struct {
   GroupID string
   URI string
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
               med.URI, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               addr, err := addr.Parse(med.URI)
               if err != nil {
                  return nil, err
               }
               med.URI = addr.String()
            }
         }
         mas.Media = append(mas.Media, med)
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
               str.Bandwidth, err = strconv.ParseInt(buf.TokenText(), 10, 64)
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
         addr, err := addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         str.URI = addr.String()
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}

func (s Stream) String() string {
   var buf []byte
   if s.Resolution != "" {
      buf = append(buf, "Resolution:"...)
      buf = append(buf, s.Resolution...)
      buf = append(buf, ' ')
   }
   buf = append(buf, "Bandwidth:"...)
   buf = strconv.AppendInt(buf, s.Bandwidth, 10)
   buf = append(buf, " Codecs:"...)
   buf = append(buf, s.Codecs...)
   if s.Audio != "" {
      buf = append(buf, " Audio:"...)
      buf = append(buf, s.Audio...)
   }
   if s.URI != "" {
      buf = append(buf, " URI:"...)
      buf = append(buf, s.URI...)
   }
   return string(buf)
}

