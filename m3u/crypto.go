package m3u

import (
   "crypto/aes"
   "crypto/cipher"
)

type Block struct {
   cipher.Block
   iv []byte
}

func NewCipher(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{block, key}, nil
}

func (b Block) Decrypt(src []byte) []byte {
   total := len(src)
   if total >= b.BlockSize() {
      dst := make([]byte, total)
      cipher.NewCBCDecrypter(b.Block, b.iv).CryptBlocks(dst, src)
      value := int(dst[total-1])
      if value < total {
         return dst[:total-value]
      }
   }
   return nil
}
