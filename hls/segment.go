package hls

import (
   "bytes"
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "io"
   "strconv"
   "strings"
   "text/scanner"
)

func (c *Cipher) ReadFrom(r io.Reader) (int64, error) {
   num, err := c.key.ReadFrom(r)
   if err != nil {
      return 0, err
   }
   c.Block, err = aes.NewCipher(c.key.Bytes())
   if err != nil {
      return 0, err
   }
   return num, nil
}

type Cipher struct {
   IV []byte
   cipher.Block
   key bytes.Buffer
}

func (s Scanner) Segment() (*Segment, error) {
   var (
      key bool
      seg Segment
   )
   for s.line.Scan() != scanner.EOF {
      line := s.line.TokenText()
      s.Init(strings.NewReader(line))
      switch {
      case strings.HasPrefix(line, "#EXT-X-KEY:"):
         key = true
         for s.Scan() != scanner.EOF {
            switch s.TokenText() {
            case "IV":
               s.Scan()
               s.Scan()
               seg.RawIV = s.TokenText()
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
         if key {
            seg.Protected = append(seg.Protected, line)
         } else {
            seg.Clear = append(seg.Clear, line)
         }
      case line == "#EXT-X-DISCONTINUITY":
         key = false
      }
   }
   return &seg, nil
}

func (s Segment) IV() ([]byte, error) {
   up := strings.ToUpper(s.RawIV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

type Segment struct {
   Clear []string
   Protected []string
   RawIV string
   RawKey string
}

func (c Cipher) Copy(w io.Writer, r io.Reader) (int, error) {
   if c.IV == nil {
      c.IV = c.key.Bytes()
   }
   buf, err := io.ReadAll(r)
   if err != nil {
      return 0, err
   }
   cipher.NewCBCDecrypter(c.Block, c.IV).CryptBlocks(buf, buf)
   if len(buf) >= 1 {
      pad := buf[len(buf)-1]
      if len(buf) >= int(pad) {
         buf = buf[:len(buf)-int(pad)]
      }
   }
   return w.Write(buf)
}
