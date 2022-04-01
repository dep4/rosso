package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "io"
   "net/url"
   "path"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

func hexDecode(s string) ([]byte, error) {
   s = strings.TrimPrefix(s, "0x")
   return hex.DecodeString(s)
}

func scanLines(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r != '\r' && r != '\n'
   }
   buf.Whitespace = 1 << '\r' | 1 << '\n'
}

func scanWords(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r == '-' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
   }
   buf.Whitespace = 1 << ' '
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

func (c Cipher) Copy(dst io.Writer, src io.Reader, iv []byte) (int, error) {
   buf, err := io.ReadAll(src)
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
   return dst.Write(buf)
}

type Information struct {
   URI *url.URL
   IV []byte
}

func (i Information) String() string {
   buf := new(strings.Builder)
   buf.WriteString("URI: ")
   buf.WriteString(i.URI.String())
   if i.IV != nil {
      buf.WriteString("\nIV: ")
      hex.NewEncoder(buf).Write(i.IV)
   }
   return buf.String()
}

type Key struct {
   Method string
   URI *url.URL
}

func (k Key) String() string {
   var buf strings.Builder
   buf.WriteString("Method: ")
   buf.WriteString(k.Method)
   buf.WriteString("\nURI: ")
   buf.WriteString(k.URI.String())
   return buf.String()
}

type Segment struct {
   Key *Key
   Info []Information
}

func NewSegment(addr *url.URL, body io.Reader) (*Segment, error) {
   var (
      buf scanner.Scanner
      err error
      info Information
      seg Segment
   )
   buf.Init(body)
   for {
      scanWords(&buf)
      if buf.Scan() == scanner.EOF {
         break
      }
      switch buf.TokenText() {
      case "EXT-X-KEY":
         seg.Key = new(Key)
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "METHOD":
               buf.Scan()
               buf.Scan()
               seg.Key.Method = buf.TokenText()
            case "URI":
               buf.Scan()
               buf.Scan()
               ref, err := strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               seg.Key.URI, err = addr.Parse(ref)
               if err != nil {
                  return nil, err
               }
            case "IV":
               buf.Scan()
               buf.Scan()
               info.IV, err = hexDecode(buf.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case "EXTINF":
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         info.URI, err = addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         seg.Info = append(seg.Info, info)
         info = Information{}
      }
   }
   return &seg, nil
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

func (s Segment) Progress(i int) (int, string) {
   pro := len(s.Info)-i
   if i == len(s.Info)-1 {
      return pro, "\n"
   }
   return pro, " "
}

func (s Segment) String() string {
   var buf strings.Builder
   if s.Key != nil {
      buf.WriteString(s.Key.String())
   }
   for _, info := range s.Info {
      buf.WriteByte('\n')
      buf.WriteString(info.String())
   }
   return buf.String()
}
