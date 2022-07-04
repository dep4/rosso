package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "io"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func (Medium) Ext() string {
   return ".m4a"
}

func (m Medium) String() string {
   var b strings.Builder
   b.WriteString("Type:")
   b.WriteString(m.Type)
   b.WriteString(" Name:")
   b.WriteString(m.Name)
   b.WriteString("\n  Group ID:")
   b.WriteString(m.Group_ID)
   if m.Characteristics != "" {
      b.WriteString("\n  Characteristics:")
      b.WriteString(m.Characteristics)
   }
   return b.String()
}

func (m Media) Get_Group_ID(value string) *Medium {
   for _, medium := range m {
      if medium.Group_ID == value {
         return &medium
      }
   }
   return nil
}

func (m Media) Get_Name(value string) *Medium {
   for _, medium := range m {
      if medium.Name == value {
         return &medium
      }
   }
   return nil
}

type Media []Medium

type Streams []Stream

func (s Streams) Bandwidth(value int) *Stream {
   distance := func(s *Stream) int {
      if s.Bandwidth > value {
         return s.Bandwidth - value
      }
      return value - s.Bandwidth
   }
   var elem *Stream
   for key, value := range s {
      if elem == nil || distance(&value) < distance(elem) {
         elem = &s[key]
      }
   }
   return elem
}

func (s Stream) String() string {
   var (
      a []string
      b string
   )
   if s.Resolution != "" {
      a = append(a, "Resolution:" + s.Resolution)
   }
   a = append(a, "Bandwidth:" + strconv.Itoa(s.Bandwidth))
   if s.Codecs != "" {
      a = append(a, "Codecs:" + s.Codecs)
   }
   if s.Audio != "" {
      b = "Audio:" + s.Audio
   }
   ja := strings.Join(a, " ")
   if b != "" {
      return ja + "\n\t" + b
   }
   return ja
}

func (Stream) Ext() string {
   return ".m4v"
}

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(s.Raw_IV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

type Segment struct {
   Key string
   Map string
   Raw_IV string
   URI []string
}

func (s Scanner) Segment() (*Segment, error) {
   var seg Segment
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      var err error
      switch {
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         seg.URI = append(seg.URI, line)
      case line == "#EXT-X-DISCONTINUITY":
         if seg.Key != "" {
            return &seg, nil
         }
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         seg.URI = nil
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               seg.Raw_IV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               seg.Key, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         s.Init(strings.NewReader(line))
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "URI":
               s.Scan()
               s.Scan()
               seg.Map, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      }
   }
   return &seg, nil
}

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
            case "CHARACTERISTICS":
               s.Scan()
               s.Scan()
               med.Characteristics, err = strconv.Unquote(s.TokenText())
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.Group_ID, err = strconv.Unquote(s.TokenText())
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
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
         s.line.Scan()
         str.URI = s.line.TokenText()
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}
