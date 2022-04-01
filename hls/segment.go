package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
   "net/url"
   "path"
   "strconv"
   "strings"
   "text/scanner"
   "unicode"
)

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

type Segment struct {
   Key *Key
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

func (i Information) String() string {
   var buf strings.Builder
   buf.WriteString("URI: ")
   buf.WriteString(i.URI.String())
   if i.IV != "" {
      buf.WriteString("\nIV: ")
      buf.WriteString(i.IV)
   }
   return buf.String()
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
               info.IV = buf.TokenText()
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
         info.IV = ""
         info.URI = nil
      }
   }
   return &seg, nil
}

type Information struct {
   URI *url.URL
   IV string
}

////////////////////////////////////////////////////////////////////////////////

type Decrypter struct {
   cipher.Block
   IV []byte
}

func NewDecrypter(src io.Reader) (*Decrypter, error) {
   key, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Decrypter{block, key}, nil
}

func (d Decrypter) Copy(dst io.Writer, src io.Reader) (int, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return 0, err
   }
   cipher.NewCBCDecrypter(d.Block, d.IV).CryptBlocks(buf, buf)
   if len(buf) >= 1 {
      pad := buf[len(buf)-1]
      if len(buf) >= int(pad) {
         buf = buf[:len(buf)-int(pad)]
      }
   }
   return dst.Write(buf)
}
