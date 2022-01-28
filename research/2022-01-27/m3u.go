package m3u

import (
   "crypto/aes"
   "crypto/cipher"
   "github.com/89z/format"
   "io"
   "net/http"
)

func unpad(buf []byte) []byte {
   total := len(buf)
   if total > 0 {
      value := int(buf[total-1])
      if value < total {
         return buf[:total-value]
      }
   }
   return nil
}

func newDecrypter(req *http.Request) (cipher.BlockMode, error) {
   logLevel.Dump(req)
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   key, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   block, err := aes.NewCipher(key)
   if err != nil {
      return nil, err
   }
   return cipher.NewCBCDecrypter(block, key), nil
}

var logLevel format.LogLevel

func writeFile(req *http.Request, dec cipher.BlockMode) ([]byte, error) {
   res, err := new(http.Transport).RoundTrip(req)
   if err != nil {
      return nil, err
   }
   defer res.Body.Close()
   src, err := io.ReadAll(res.Body)
   if err != nil {
      return nil, err
   }
   dst := make([]byte, len(src))
   dec.CryptBlocks(dst, src)
   return unpad(dst), nil
}
