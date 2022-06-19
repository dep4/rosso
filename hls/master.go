package hls

import (
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

type Medium struct {
   Type string
   Name string
   Group_ID string
   Raw_URI string
}

type Stream struct {
   Resolution string
   Bandwidth int64 // handle duplicate resolution
   Video_Range string // handle duplicate bandwidth
   Raw_Codecs string // handle missing resolution
   Raw_URI string
}

const (
   AAC = ".aac"
   TS = ".ts"
)

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

func (m Media) Get_Group_ID(val string) *Medium {
   for _, medium := range m {
      if medium.Group_ID == val {
         return &medium
      }
   }
   return nil
}

// English
func (m Media) Get_Name(val string) *Medium {
   for _, medium := range m {
      if medium.Name == val {
         return &medium
      }
   }
   return nil
}

// stereo
func (m Media) Group_ID(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.Group_ID, val) {
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
func (m Media) URI(val string) Media {
   var out Media
   for _, medium := range m {
      if strings.Contains(medium.Raw_URI, val) {
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

func (m Medium) String() string {
   var buf strings.Builder
   buf.WriteString("Type:")
   buf.WriteString(m.Type)
   buf.WriteString(" Name:")
   buf.WriteString(m.Name)
   buf.WriteString(" ID:")
   buf.WriteString(m.Group_ID)
   return buf.String()
}

type Scanner struct {
   line scanner.Scanner
   scanner.Scanner
}

func New_Scanner(body io.Reader) Scanner {
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
               med.Group_ID, err = strconv.Unquote(s.TokenText())
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
               med.Raw_URI, err = strconv.Unquote(s.TokenText())
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
               str.Raw_Codecs, err = strconv.Unquote(s.TokenText())
            case "VIDEO-RANGE":
               s.Scan()
               s.Scan()
               str.Video_Range = s.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         s.line.Scan()
         str.Raw_URI = s.line.TokenText()
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}

func (s Stream) Codecs() string {
   codecs := strings.Split(s.Raw_Codecs, ",")
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
   buf = append(buf, s.Video_Range...)
   if s.Raw_Codecs != "" {
      buf = append(buf, " Codecs:"...)
      buf = append(buf, s.Codecs()...)
   }
   return string(buf)
}

type Streams []Stream

// hvc1 mp4a
func (s Streams) Codecs(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Raw_Codecs, val) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) Get_Bandwidth(val int64) *Stream {
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
func (s Streams) URI(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Raw_URI, val) {
         out = append(out, stream)
      }
   }
   return out
}

// PQ
func (s Streams) Video_Range(val string) Streams {
   var out Streams
   for _, stream := range s {
      if stream.Video_Range == val {
         out = append(out, stream)
      }
   }
   return out
}
