package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "encoding/hex"
   "io"
   "strconv"
   "strings"
   "text/scanner"
)

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

type Block struct {
   cipher.Block
   key []byte
}

func NewBlock(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{block, key}, nil
}

func (b Block) Mode(r io.Reader, iv []byte) *BlockMode {
   var mode BlockMode
   mode.BlockMode = cipher.NewCBCDecrypter(b.Block, iv)
   mode.reader = r
   return &mode
}

func (b Block) ModeKey(r io.Reader) *BlockMode {
   return b.Mode(r, b.key)
}

type BlockMode struct {
   cipher.BlockMode
   clear []byte
   protected []byte
   reader io.Reader
}

func (b *BlockMode) Read(p []byte) (int, error) {
   // move to protected
   num, err := b.reader.Read(p)
   b.protected = append(b.protected, p[:num]...)
   // move to clear
   num = len(b.protected) - len(b.protected) % 16
   b.CryptBlocks(b.protected, b.protected[:num])
   b.clear = append(b.clear, b.protected[:num]...)
   b.protected = b.protected[num:]
   // move to out
   num = copy(p, b.clear)
   b.clear = b.clear[num:]
   return num, err
}
