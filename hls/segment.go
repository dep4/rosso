package hls

import (
   "crypto/aes"
   "crypto/cipher"
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

func (s Segment) Key(base *url.URL) (*url.URL, error) {
   return base.Parse(*s.RawKey)
}

func (s Segment) URI(base *url.URL) (*url.URL, error) {
   return base.Parse(s.RawURI)
}

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(*s.RawIV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

func (s Scanner) Segments() (Segments, error) {
   var (
      seg Segment
      segs Segments
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
               text := s.TokenText()
               seg.RawIV = &text
            case "URI":
               s.Scan()
               s.Scan()
               text, err := strconv.Unquote(s.TokenText())
               if err != nil {
                  return nil, err
               }
               seg.RawKey = &text
            }
         }
      case len(line) >= 1 && !strings.HasPrefix(line, "#"):
         seg.RawURI = line
         segs = append(segs, seg)
      case line == "#EXT-X-DISCONTINUITY":
         seg.RawIV = nil
         seg.RawKey = nil
      }
   }
   return segs, nil
}

type Segments []Segment

func (s Segments) Key() Segments {
   var segs Segments
   for _, seg := range s {
      if seg.RawKey != nil {
         segs = append(segs, seg)
      }
   }
   return segs
}

type Segment struct {
   RawURI string
   RawIV *string
   RawKey *string
}
