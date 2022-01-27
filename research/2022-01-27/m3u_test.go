package m3u

import (
   "crypto/aes"
   "crypto/cipher"
   "os"
   "testing"
)

func TestDecrypt(t *testing.T) {
   src, err := segment()
   if err != nil {
      t.Fatal(err)
   }
   key, err := cryptKey()
   if err != nil {
      t.Fatal(err)
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      t.Fatal(err)
   }
   dst := make([]byte, len(src))
   cipher.NewCBCDecrypter(block, key).CryptBlocks(dst, src)
   dst = unpad(dst)
   if err := os.WriteFile("segment1_1_av.ts", dst, os.ModePerm); err != nil {
      t.Fatal(err)
   }
}
