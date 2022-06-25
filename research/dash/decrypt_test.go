package dash

import (
   "encoding/hex"
   "os"
   "testing"
)

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
   key, err := hex.DecodeString("680a46ebd6cf2b9a6a0b05a24dcf944a")
   if err != nil {
      t.Fatal(err)
   }
   if err := decryptMP4withCenc(enc, key, dec); err != nil {
      t.Fatal(err)
   }
}
