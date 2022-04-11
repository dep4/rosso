package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
   "net/url"
   "path"
   "strconv"
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

func (c Cipher) Decrypt(info Information, src io.Reader) ([]byte, error) {
   buf, err := io.ReadAll(src)
   if err != nil {
      return nil, err
   }
   if info.IV == nil {
      info.IV = c.key
   }
   cipher.NewCBCDecrypter(c.Block, info.IV).CryptBlocks(buf, buf)
   if len(buf) >= 1 {
      pad := buf[len(buf)-1]
      if len(buf) >= int(pad) {
         buf = buf[:len(buf)-int(pad)]
      }
   }
   return buf, nil
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

func (s *Scanner) Segment(addr *url.URL) (*Segment, error) {
   var (
      err error
      info Information
      seg Segment
   )
   for {
      s.splitWords()
      if s.Scan() == scanner.EOF {
         break
      }
      switch s.TokenText() {
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
      case "EXT-X-KEY":
         for s.Scan() != '\n' {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               info.IV, err = hexDecode(s.TokenText())
               if err != nil {
                  return nil, err
               }
            case "URI":
               s.Scan()
               s.Scan()
               ref, err := strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
               seg.Key, err = addr.Parse(ref)
               if err != nil {
                  return nil, err
               }
            }
         }
      }
   }
   return &seg, nil
}
