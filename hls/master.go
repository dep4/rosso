package hls

import (
   "encoding/hex"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func scanHex(s string) ([]byte, error) {
   up := strings.ToUpper(s)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func scanURL(s string, addr *url.URL) (*url.URL, error) {
   ref, err := strconv.Unquote(s)
   if err != nil {
      return nil, err
   }
   return addr.Parse(ref)
}

func (s *Scanner) splitLines() {
   s.IsIdentRune = func(r rune, i int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      return true
   }
   s.Whitespace |= 1 << '\n'
   s.Whitespace |= 1 << '\r'
}

func (s *Scanner) splitWords() {
   s.IsIdentRune = func(r rune, i int) bool {
      if r == '-' {
         return true
      }
      if r == '.' {
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
   s.Whitespace = 1 << ' '
}

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

func (m Media) GetGroupID(val string) *Medium {
   for _, medium := range m {
      if medium.GroupID == val {
         return &medium
      }
   }
   return nil
}

// stereo
func (m Media) GroupID(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.GroupID, val) {
         out = append(out, medium)
      }
   }
   return out
}

// English
func (m Media) Name(val string) Media {
   var out Media
   for _, medium := range m {
      if medium.Name == val {
         out = append(out, medium)
      }
   }
   return out
}

// cdn
func (m Media) RawQuery(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.URI.RawQuery, val) {
         out = append(out, medium)
      }
   }
   return out
}

// AUDIO
func (m Media) Type(val string) Media {
   var out Media
   for _, medium := range m {
      if medium.Type == val {
         out = append(out, medium)
      }
   }
   return out
}

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
   scanner.Scanner
}

func NewScanner(body io.Reader) *Scanner {
   var scan Scanner
   scan.Init(body)
   return &scan
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
         var med Medium
         for s.Scan() != '\n' {
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
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
            case "VIDEO-RANGE":
               s.Scan()
               s.Scan()
               str.VideoRange = s.TokenText()
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
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
      fmt.Fprint(f, " URI:", s.URI)
   }
}

type Streams []Stream

// hvc1 mp4a
func (s Streams) Codecs(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Codecs, val) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) GetBandwidth(val int64) *Stream {
   distance := func(s *Stream) int64 {
      if s.Bandwidth > val {
         return s.Bandwidth - val
      }
      return val - s.Bandwidth
   }
   var out *Stream
   for key, val := range s {
      if out == nil || distance(&val) < distance(out) {
         out = &s[key]
      }
   }
   return out
}

// cdn=vod-ak-aoc.tv.apple.com
func (s Streams) RawQuery(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.URI.RawQuery, val) {
         out = append(out, stream)
      }
   }
   return out
}

// PQ
func (s Streams) VideoRange(val string) Streams {
   var out Streams
   for _, stream := range s {
      if stream.VideoRange == val {
         out = append(out, stream)
      }
   }
   return out
}
