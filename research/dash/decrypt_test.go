package decrypt

import (
   "encoding/hex"
   "os"
   "testing"
)

const rawKey = "6b1f79ba70956a37fe716997b8d211ae"

func TestDecrypt(t *testing.T) {
   enc, err := os.Open("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer enc.Close()
   dec, err := os.Create("ignore/dec.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dec.Close()
   key, err := hex.DecodeString(rawKey)
   if err != nil {
      t.Fatal(err)
   }
   if err := decryptMP4withCenc(enc, key, dec); err != nil {
      t.Fatal(err)
   }
}
