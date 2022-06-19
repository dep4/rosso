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
               seg.Raw_IV = s.TokenText()
            case "URI":
               s.Scan()
               s.Scan()
               var err error
               seg.Raw_Key, err = strconv.Unquote(s.TokenText())
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
   up := strings.ToUpper(s.Raw_IV)
   return hex.DecodeString(strings.TrimPrefix(up, "0X"))
}

type Segment struct {
   Clear []string
   Protected []string
   Raw_IV string
   Raw_Key string
}

type Block struct {
   cipher.Block
   key []byte
}

func New_Block(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{block, key}, nil
}

func (b Block) Mode(r io.Reader, iv []byte) *Block_Mode {
   var mode Block_Mode
   mode.BlockMode = cipher.NewCBCDecrypter(b.Block, iv)
   mode.reader = r
   return &mode
}

func (b Block) Mode_Key(r io.Reader) *Block_Mode {
   return b.Mode(r, b.key)
}

type Block_Mode struct {
   cipher.BlockMode
   reader io.Reader
   cipher []byte
   plain []byte
}

func (b Block_Mode) len_message(err error) int {
   num := len(b.plain)
   pad := b.plain[num-1]
   if err == nil {
      pad = 16
   }
   return num - int(pad)
}

func (b Block_Mode) len_plain() int {
   num := len(b.cipher)
   return num - num % 16
}

func (b *Block_Mode) Read(p []byte) (int, error) {
   // ciphertext length
   num, err := b.reader.Read(p)
   b.cipher = append(b.cipher, p[:num]...)
   // plaintext length
   num = b.len_plain()
   b.CryptBlocks(b.cipher, b.cipher[:num])
   b.plain = append(b.plain, b.cipher[:num]...)
   b.cipher = b.cipher[num:]
   // message length
   num = b.len_message(err)
   // output length
   num = copy(p, b.plain[:num])
   b.plain = b.plain[num:]
   return num, err
}
