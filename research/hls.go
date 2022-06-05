package hls

import (
   "bufio"
   "bytes"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "text/scanner"
   "unicode"
)

func isMedia(s []byte) bool {
   prefix := []byte("#EXT-X-MEDIA:")
   return bytes.HasPrefix(s, prefix)
}

func isStream(s []byte) bool {
   prefix := []byte("#EXT-X-STREAM-INF:")
   return bytes.HasPrefix(s, prefix)
}

func scanURL(base *url.URL, raw string) (*url.URL, error) {
   ref, err := strconv.Unquote(raw)
   if err != nil {
      return nil, err
   }
   return base.Parse(ref)
}

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

type Medium struct {
   Type string
   Name string
   GroupID string
   URI *url.URL
}

func (m Medium) Format(f fmt.State, verb rune) {
   fmt.Fprint(f, "Type:", m.Type)
   fmt.Fprint(f, " Name:", m.Name)
   fmt.Fprint(f, " ID:", m.GroupID)
   if verb == 'a' {
      fmt.Fprint(f, " URI:", m.URI)
   }
}

type Scanner struct {
   bufio *bufio.Scanner
   scanner.Scanner
}

func NewScanner(body io.Reader) Scanner {
   var scan Scanner
   scan.bufio = bufio.NewScanner(body)
   scan.IsIdentRune = func(r rune, i int) bool {
      if r == '-' {
         return true
      }
      if unicode.IsDigit(r) {
         return true
      }
      if unicode.IsLetter(r) {
         return true
      }
      return false
   }
   return scan
}

func (s Scanner) Master(base *url.URL) (*Master, error) {
   var (
      err error
      mas Master
   )
   for s.bufio.Scan() {
      slice := s.bufio.Bytes()
      s.Init(bytes.NewReader(slice))
      switch {
      case isMedia(slice):
         var med Medium
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.GroupID, err = strconv.Unquote(s.TokenText())
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
            case "URI":
               s.Scan()
               s.Scan()
               med.URI, err = scanURL(base, s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case isStream(slice):
         var str Stream
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            case "VIDEO-RANGE":
               s.Scan()
               s.Scan()
               str.VideoRange = s.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         s.bufio.Scan()
         str.URI, err = base.Parse(s.bufio.Text())
         if err != nil {
            return nil, err
         }
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}

type Stream struct {
   Resolution string
   VideoRange string // handle duplicate bandwidth
   Bandwidth int64 // handle duplicate resolution
   Codecs string // handle missing resolution
   URI *url.URL
}

func (s Stream) Format(f fmt.State, verb rune) {
   if s.Resolution != "" {
      fmt.Fprint(f, "Resolution:", s.Resolution, " ")
   }
   fmt.Fprint(f, "Bandwidth:", s.Bandwidth)
   if s.Codecs != "" {
      fmt.Fprint(f, " Codecs:", s.Codecs)
   }
   if verb == 'a' {
      fmt.Fprint(f, " Range:", s.VideoRange)
   }
   if verb == 'b' {
      fmt.Fprint(f, " URI:", s.URI)
   }
}

type Streams []Stream

type Segment struct {
   Key *url.URL
   Info []Information
}

type Information struct {
   IV []byte
   URI *url.URL
}

func (s *Scanner) Segment(addr *url.URL) (*Segment, error) {
   var (
      info Information
      seg Segment
   )
   for {
      s.splitWords()
      if s.Scan() == scanner.EOF {
         break
      }
      var err error
      switch s.TokenText() {
      case "EXT-X-KEY":
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               info.IV, err = scanHex(s.TokenText())
            case "URI":
               s.Scan()
               s.Scan()
               seg.Key, err = scanURL(s.TokenText(), addr)
            }
            if err != nil {
               return nil, err
            }
         }
      case "EXTINF":
         s.splitLines()
         s.Scan()
         s.Scan()
         info.URI, err = addr.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
         seg.Info = append(seg.Info, info)
         info = Information{}
      }
   }
   return &seg, nil
}
