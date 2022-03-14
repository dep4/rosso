package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
   "net/url"
   "path"
   "strconv"
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

type Information struct {
   Duration string
   URI *url.URL
}

type Key struct {
   Method string
   URI *url.URL
}

type Segment struct {
   Key *Key
   Info []Information
}

func NewSegment(addr *url.URL, body io.Reader) (*Segment, error) {
   var (
      buf scanner.Scanner
      err error
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
            }
         }
      case "EXTINF":
         var info Information
         buf.Scan()
         buf.Scan()
         info.Duration = buf.TokenText()
         scanLines(&buf)
         buf.Scan()
         buf.Scan()
         info.URI, err = addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         seg.Info = append(seg.Info, info)
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
