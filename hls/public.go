package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "fmt"
   "io"
   "net/url"
   "strconv"
   "text/scanner"
   "time"
)

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

type Segment struct {
   Key *url.URL
   Info []Information
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

type Master struct {
   Stream []Stream
   Media []Media
}

func (m Master) Audio(str Stream) *Media {
   for _, med := range m.Media {
      if med.GroupID == str.Audio {
         return &med
      }
   }
   return nil
}

func (m Master) Len() int {
   return len(m.Stream)
}

func (m Master) Swap(i, j int) {
   m.Stream[i], m.Stream[j] = m.Stream[j], m.Stream[i]
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
            case "GROUP-ID":
               s.Scan()
               s.Scan()
               med.GroupID, err = strconv.Unquote(s.TokenText())
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
         s.splitLines()
         s.Scan()
         str.URI, err = addr.Parse(s.TokenText())
         if err != nil {
            return nil, err
         }
         mas.Stream = append(mas.Stream, str)
      }
   }
   return &mas, nil
}
