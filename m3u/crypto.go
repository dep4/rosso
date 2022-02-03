package m3u

import (
   "crypto/aes"
   "crypto/cipher"
)

type Block struct {
   IV []byte
   cipher.Block
}

func NewCipher(key []byte) (*Block, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return &Block{key, block}, nil
}

// We do not care about the ciphertext, so this works in place.
func (b Block) Decrypt(src []byte) []byte {
   total := len(src)
   if total >= b.BlockSize() {
      cipher.NewCBCDecrypter(b.Block, b.IV).CryptBlocks(src, src)
      value := int(src[total-1])
      if value < total {
         return src[:total-value]
      }
   }
   return nil
}
