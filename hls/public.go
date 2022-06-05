package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "fmt"
   "io"
   "net/url"
   "strconv"
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
   Streams []Stream
   Media []Media
}

type Media struct {
   Name string
   Type string
   URI *url.URL
}

func (m Media) Format(f fmt.State, verb rune) {
   fmt.Fprint(f, "Type:", m.Type)
   fmt.Fprint(f, " Name:", m.Name)
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
         var med Media
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "TYPE":
               s.Scan()
               s.Scan()
               med.Type = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               med.URI, err = scanURL(s.TokenText(), addr)
            case "NAME":
               s.Scan()
               s.Scan()
               med.Name, err = strconv.Unquote(s.TokenText())
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
            case "CODECS":
               s.Scan()
               s.Scan()
               str.Codecs, err = strconv.Unquote(s.TokenText())
            case "BANDWIDTH":
               s.Scan()
               s.Scan()
               str.Bandwidth, err = strconv.ParseInt(s.TokenText(), 10, 64)
            case "RESOLUTION":
               s.Scan()
               s.Scan()
               str.Resolution = s.TokenText()
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
      fmt.Fprint(f, " URI:", s.URI)
   }
}

// Provide a name such as "English"
func (m Master) Audio(name string) *Media {
   for _, med := range m.Media {
      if med.Type == "AUDIO" && med.Name == name {
         return &med
      }
   }
   return nil
}

func (m Master) Stream(bandwidth int64) *Stream {
   distance := func(s *Stream) int64 {
      if s.Bandwidth > bandwidth {
         return s.Bandwidth - bandwidth
      }
      return bandwidth - s.Bandwidth
   }
   var dst *Stream
   for i, src := range m.Streams {
      if dst == nil || distance(&src) < distance(dst) {
         dst = &m.Streams[i]
      }
   }
   return dst
}
