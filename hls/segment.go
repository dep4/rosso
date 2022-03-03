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

func (s Segment) Ext() (string, error) {
   for _, info := range s.Info {
      addr, err := url.Parse(info.URI)
      if err != nil {
         return "", err
      }
      ext := path.Ext(addr.Path)
      if ext != "" {
         return ext, nil
      }
   }
   return "", notPresent{"path.Ext"}
}

type notPresent struct {
   value string
}

func (n notPresent) Error() string {
   return strconv.Quote(n.value) + " is not present"
}

type Segment struct {
   Key struct {
      Method string
      URI string
   }
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
         for buf.Scan() != '\n' {
            switch buf.TokenText() {
            case "METHOD":
               buf.Scan()
               buf.Scan()
               seg.Key.Method = buf.TokenText()
            case "URI":
               buf.Scan()
               buf.Scan()
               seg.Key.URI, err = strconv.Unquote(buf.TokenText())
               if err != nil {
                  return nil, err
               }
               addr, err := addr.Parse(seg.Key.URI)
               if err != nil {
                  return nil, err
               }
               seg.Key.URI = addr.String()
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
         addr, err := addr.Parse(buf.TokenText())
         if err != nil {
            return nil, err
         }
         info.URI = addr.String()
         seg.Info = append(seg.Info, info)
      }
   }
   return &seg, nil
}

type Information struct {
   Duration string
   URI string
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

type Decrypter struct {
   cipher.Block
   IV []byte
}

func scanLines(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r != '\n'
   }
   buf.Whitespace = 1 << '\n'
}

func scanWords(buf *scanner.Scanner) {
   buf.IsIdentRune = func(r rune, i int) bool {
      return r == '-' || r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r)
   }
   buf.Whitespace = 1 << ' '
}