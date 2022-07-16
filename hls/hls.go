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

func (self Block) Decrypt(text, iv []byte) []byte {
   cipher.NewCBCDecrypter(self.Block, iv).CryptBlocks(text, text)
   if len(text) >= 1 {
      pad := text[len(text)-1]
      if len(text) >= int(pad) {
         text = text[:len(text)-int(pad)]
      }
   }
   return text
}

func (self Block) Decrypt_Key(text []byte) []byte {
   return self.Decrypt(text, self.key)
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

func (self Scanner) Master() (*Master, error) {
   var mas Master
   for self.line.Scan() != scanner.EOF {
      var err error
      line := self.line.TokenText()
      self.Init(strings.NewReader(line))
      switch {
      case strings.HasPrefix(line, "#EXT-X-MEDIA:"):
         var med Medium
         for self.Scan() != scanner.EOF {
            switch self.TokenText() {
            case "CHARACTERISTICS":
               self.Scan()
               self.Scan()
               med.Characteristics, err = strconv.Unquote(self.TokenText())
            case "GROUP-ID":
               self.Scan()
               self.Scan()
               med.Group_ID, err = strconv.Unquote(self.TokenText())
            case "NAME":
               self.Scan()
               self.Scan()
               med.Name, err = strconv.Unquote(self.TokenText())
            case "TYPE":
               self.Scan()
               self.Scan()
               med.Type = self.TokenText()
            case "URI":
               self.Scan()
               self.Scan()
               med.Raw_URI, err = strconv.Unquote(self.TokenText())
            }
            if err != nil {
               return nil, err
            }
         }
         mas.Media = append(mas.Media, med)
      case strings.HasPrefix(line, "#EXT-X-STREAM-INF:"):
         var str Stream
         for self.Scan() != scanner.EOF {
            switch self.TokenText() {
            case "AUDIO":
               self.Scan()
               self.Scan()
               str.Audio, err = strconv.Unquote(self.TokenText())
            case "BANDWIDTH":
               self.Scan()
               self.Scan()
               str.Bandwidth, err = strconv.Atoi(self.TokenText())
            case "CODECS":
               self.Scan()
               self.Scan()
               str.Codecs, err = strconv.Unquote(self.TokenText())
            case "RESOLUTION":
               self.Scan()
               self.Scan()
               str.Resolution = self.TokenText()
            }
            if err != nil {
               return nil, err
            }
         }
         self.line.Scan()
         str.Raw_URI = self.line.TokenText()
         mas.Streams = append(mas.Streams, str)
      }
   }
   return &mas, nil
}

func (self Scanner) Segment() (*Segment, error) {
   var seg Segment
   for self.line.Scan() != scanner.EOF {
      line := self.line.TokenText()
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
         self.Init(strings.NewReader(line))
         for self.Scan() != scanner.EOF {
            switch self.TokenText() {
            case "IV":
               self.Scan()
               self.Scan()
               seg.Raw_IV = self.TokenText()
            case "URI":
               self.Scan()
               self.Scan()
               seg.Key, err = strconv.Unquote(self.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case strings.HasPrefix(line, "#EXT-X-MAP:"):
         self.Init(strings.NewReader(line))
         for self.Scan() != scanner.EOF {
            switch self.TokenText() {
            case "URI":
               self.Scan()
               self.Scan()
               seg.Map, err = strconv.Unquote(self.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      }
   }
   return &seg, nil
}

type Segment struct {
   Key string
   Map string
   Raw_IV string
   URI []string
}

func (self Segment) IV() ([]byte, error) {
   up := strings.ToUpper(self.Raw_IV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}
