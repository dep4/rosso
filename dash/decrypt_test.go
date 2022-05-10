package dash

import (
   "encoding/hex"
   "os"
   "testing"
)

const hexKey = "13d7c7cf295444944b627ef0ad2c1b3c"

func TestDecrypt(t *testing.T) {
   src, err := os.Open("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer src.Close()
   dst, err := os.Create("ignore/dec.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dst.Close()
   key, err := hex.DecodeString(hexKey)
   if err != nil {
      t.Fatal(err)
   }
   if err := Decrypt(dst, src, key); err != nil {
      t.Fatal(err)
   }
}
