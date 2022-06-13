package hls

import (
   "crypto/aes"
   "crypto/cipher"
   "io"
)

type Cipher struct {
   cipher.BlockMode
   reader io.Reader
}

func newCipher(r io.Reader, key []byte) (*Cipher, error) {
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   var c Cipher
   c.BlockMode = cipher.NewCBCDecrypter(block, key)
   c.reader = r
   return &c, nil
}

func (c Cipher) Read(p []byte) (int, error) {
   var high int
   for high < len(p) {
      high += 16
   }
   n, err := c.reader.Read(p[:high])
   c.CryptBlocks(p, p[:high])
   return n, err
}
