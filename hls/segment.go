package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/base64"
   "encoding/hex"
   "io"
   "net/url"
   "strconv"
   "strings"
   "text/scanner"
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

func (s Scanner) Segment() (*Segment, error) {
   var (
      info Information
      seg Segment
   )
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      s.Init(strings.NewReader(line))
      switch {
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               info.RawIV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               var err error
               seg.RawKey, err = strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
            }
         }
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         info.RawURI = line
         seg.Info = append(seg.Info, info)
         info = Information{}
      }
   }
   return &seg, nil
}

type Information struct {
   RawIV string
   RawURI string
}

func (i Information) IV() ([]byte, error) {
   up := strings.ToUpper(i.RawIV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func (i Information) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(i.RawURI)
}

type Segment struct {
   Info []Information
   RawKey string
}

func (s Segment) Key(base *url.URL) (*url.URL, error) {
   return base.Parse(s.RawKey)
}

func (s Segment) PSSH() ([]byte, error) {
   _, after, found := strings.Cut(s.RawKey, "data:text/plain;base64,")
   if found {
      s.RawKey = after
   }
   return base64.StdEncoding.DecodeString(s.RawKey)
}
