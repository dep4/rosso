package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
)

const (
   AAC = ".aac"
   TS = ".ts"
)

type Cipher struct {
   cipher.Block
   key []byte
}

func NewCipher(src io.Reader) (*Cipher, error) {
   key, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Cipher{block, key}, nil
}

func (c Cipher) Copy(w io.Writer, r io.Reader, iv []byte) (int, error) {
   buf, err := io.ReadAll(r)
   if err != nil {
      return 0, err
   }
   if iv == nil {
      iv = c.key
   }
   cipher.NewCBCDecrypter(c.Block, iv).CryptBlocks(buf, buf)
   if len(buf) >= 1 {
      pad := buf[len(buf)-1]
      if len(buf) >= int(pad) {
         buf = buf[:len(buf)-int(pad)]
      }
   }
   return w.Write(buf)
}

type Information struct {
   IV []byte
   URI *url.URL
}

type Master struct {
   Media Media
   Streams Streams
}

type Media []Medium

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

func (m Media) Medium(groupID string) *Medium {
   for _, medium := range m {
      if medium.GroupID == groupID {
         return &medium
      }
   }
   return nil
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

type Segment struct {
   Key *url.URL
   Info []Information
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
func (s Streams) Codec(val string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Codecs, val) {
         out = append(out, stream)
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

func (s Streams) Stream(bandwidth int64) *Stream {
   distance := func(s *Stream) int64 {
      if s.Bandwidth > bandwidth {
         return s.Bandwidth - bandwidth
      }
      return bandwidth - s.Bandwidth
   }
   var out *Stream
   for key, val := range s {
      if out == nil || distance(&val) < distance(out) {
         out = &s[key]
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
