package hls

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

type Block struct {
   cipher.Block
   key []byte
}

func New_Block(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{block, key}, nil
}

func (b Block) Decrypt(text, iv []byte) []byte {
   cipher.NewCBCDecrypter(b.Block, iv).CryptBlocks(text, text)
   if len(text) >= 1 {
      pad := text[len(text)-1]
      if len(text) >= int(pad) {
         text = text[:len(text)-int(pad)]
      }
   }
   return text
}

func (b Block) Decrypt_Key(text []byte) []byte {
   return b.Decrypt(text, b.key)
}

func (s Streams) Audio(value string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Audio, value) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) Codecs(value string) Streams {
   var out Streams
   for _, stream := range s {
      if strings.Contains(stream.Codecs, value) {
         out = append(out, stream)
      }
   }
   return out
}

func (s Streams) Get_Bandwidth(value int64) *Stream {
   distance := func(s *Stream) int64 {
      if s.Bandwidth > value {
         return s.Bandwidth - value
      }
      return value - s.Bandwidth
   }
   var out *Stream
   for key, value := range s {
      if out == nil || distance(&value) < distance(out) {
         out = &s[key]
      }
   }
   return out
}

type Stream struct {
   Audio string
   Bandwidth int64
   Codecs string
   Resolution string
   URI string
}

type Streams []Stream

func (s Stream) String() string {
   var b []byte
   b = append(b, "Bandwidth:"...)
   b = strconv.AppendInt(b, s.Bandwidth, 10)
   if s.Codecs != "" {
      b = append(b, " Codecs:"...)
      b = append(b, s.Codecs...)
   }
   if s.Resolution != "" {
      b = append(b, " Resolution:"...)
      b = append(b, s.Resolution...)
   }
   b = append(b, "\n\tAudio:"...)
   b = append(b, s.Audio...)
   return string(b)
}

func (s Stream) Ext(b []byte) string {
   if bytes.Contains(b, []byte("ftypiso5")) {
      return ".m4v"
   }
   if bytes.HasPrefix(b, []byte{'G'}) {
      return ".ts"
   }
   return ""
}

type Master struct {
   Media Media
   Streams Streams
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
               med.URI, err = strconv.Unquote(s.TokenText())
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
            case "AUDIO":
               s.Scan()
               s.Scan()
               str.Audio, err = strconv.Unquote(s.TokenText())
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
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
         s.line.Scan()
         str.URI = s.line.TokenText()
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}
