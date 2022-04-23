package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "fmt"
   "io"
   "net/url"
   "path"
   "strconv"
   "strings"
   "text/scanner"
   "time"
   "unicode"
)

func scanHex(s string) ([]byte, error) {
   up := strings.ToUpper(s)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func scanDuration(s string) (time.Duration, error) {
   sec, err := strconv.ParseFloat(s, 64)
   if err != nil {
      return 0, err
   }
   return time.Duration(sec * 1000) * time.Millisecond, nil
}

func scanURL(s string, addr *url.URL) (*url.URL, error) {
   ref, err := strconv.Unquote(s)
   if err != nil {
      return nil, err
   }
   return addr.Parse(ref)
}

type Bandwidth struct {
   *Master
   Target int
}

func (b Bandwidth) Less(i, j int) bool {
   distance := func(k int) int {
      diff := b.Stream[k].Bandwidth - b.Target
      if diff >= 0 {
         return diff
      }
      return -diff
   }
   return distance(i) < distance(j)
}

type Information struct {
   IV []byte
   // If we embed this, it will hijack String method
   Duration time.Duration
   URI *url.URL
}

type Media struct {
   GroupID string
   URI *url.URL
}

type Scanner struct {
   scanner.Scanner
}

func NewScanner(body io.Reader) *Scanner {
   var scan Scanner
   scan.Init(body)
   return &scan
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
         s.Scan()
         s.Scan()
         info.Duration, err = scanDuration(s.TokenText())
         if err != nil {
            return nil, err
         }
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

type Segment struct {
   Key *url.URL
   Info []Information
}

func (s Segment) Ext() string {
   for _, info := range s.Info {
      ext := path.Ext(info.URI.Path)
      if ext != "" {
         return ext
      }
   }
   return ""
}

func (s Segment) Length(str Stream) int64 {
   var dur time.Duration
   for _, info := range s.Info {
      dur += info.Duration
   }
   length := float64(str.Bandwidth) / 8 * dur.Seconds()
   return int64(length)
}

type Stream struct {
   Resolution string
   Bandwidth int // handle duplicate resolution
   Codecs string // handle missing resolution
   Audio string // link to Media
   URI *url.URL
}

func (s Stream) Format(f fmt.State, verb rune) {
   if s.Resolution != "" {
      fmt.Fprint(f, "Resolution:", s.Resolution, " ")
   }
   fmt.Fprint(f, "Bandwidth:", s.Bandwidth)
   fmt.Fprint(f, " Codecs:", s.Codecs)
   if verb == 'a' {
      fmt.Fprint(f, " Audio:", s.Audio)
      fmt.Fprint(f, " URI:", s.URI)
   }
}

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
