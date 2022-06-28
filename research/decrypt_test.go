package dash

import (
   "encoding/hex"
   "os"
   "testing"
)

const raw_key = "22bdb0063805260307ee5045c0f3835a"

func Test_Decrypt(t *testing.T) {
   enc, err := os.Open("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer enc.Close()
   dec, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dec.Close()
   key, err := hex.DecodeString(raw_key)
   if err != nil {
      t.Fatal(err)
   }
   if err := decryptMP4withCenc(enc, key, dec); err != nil {
      t.Fatal(err)
   }
}
