package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
)

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
