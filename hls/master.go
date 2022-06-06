package hls

import (
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

const (
   AAC = ".aac"
   TS = ".ts"
)

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
func (m Media) RawURI(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.RawURI, val) {
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
   RawURI string
}

func (m Medium) String() string {
   var buf strings.Builder
   buf.WriteString("Type:")
   buf.WriteString(m.Type)
   buf.WriteString(" Name:")
   buf.WriteString(m.Name)
   buf.WriteString(" ID:")
   buf.WriteString(m.GroupID)
   return buf.String()
}

func (m Medium) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(m.RawURI)
}

type Scanner struct {
   line scanner.Scanner
   scanner.Scanner
}

func NewScanner(body io.Reader) Scanner {
   var scan Scanner
   scan.line.Init(body)
   scan.line.IsIdentRune = func(r rune, i int) bool {
      if r == '\n' {
         return false
      }
      if r == '\r' {
         return false
      }
      if r == scanner.EOF {
         return false
      }
      return true
   }
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

func (s Scanner) Master() (*Master, error) {
   var mas Master
   for s.line.Scan() != scanner.EOF {
      var err error
      line := s.line.TokenText()
      s.Init(strings.NewReader(line))
      switch {
      case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
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
               med.RawURI, err = strconv.Unquote(s.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
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
               str.RawCodecs, err = strconv.Unquote(s.TokenText())
            case "VIDEO-RANGE":
               s.Scan()
               s.Scan()
               str.VideoRange = s.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         s.line.Scan()
         str.RawURI = s.line.TokenText()
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   VideoRange string // handle duplicate bandwidth
   RawCodecs string // handle missing resolution
   RawURI string
}

func (s Stream) Codecs() string {
   codecs := strings.Split(s.RawCodecs, ",")
   for i, codec := range codecs {
      before, _, found := strings.Cut(codec, ".")
      if found {
         codecs[i] = before
      }
   }
   return strings.Join(codecs, ",")
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
   buf = append(buf, " Range:"...)
   buf = append(buf, s.VideoRange...)
   if s.RawCodecs != "" {
      buf = append(buf, " Codecs:"...)
      buf = append(buf, s.Codecs()...)
   }
   return string(buf)
}

func (s Stream) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(s.RawURI)
}

type Streams []Stream

// hvc1 mp4a
func (s Streams) Codecs(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.RawCodecs, val) {
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
func (s Streams) RawURI(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.RawURI, val) {
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
