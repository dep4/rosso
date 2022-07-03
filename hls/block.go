package hls

import (
   "crypto/aes"
   "crypto/cipher"
)

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

func (b Block) Decrypt(text, iv []byte) []byte {
   cipher.NewCBCDecrypter(b.Block, iv).CryptBlocks(text, text)
   if len(text) >= 1 {
      pad := text[len(text)-1]
      if len(text) >= int(pad) {
         text = text[:len(text)-int(pad)]
      }
   }
   return text
}

func (b Block) Decrypt_Key(text []byte) []byte {
   return b.Decrypt(text, b.key)
}
