package dash

import (
   "encoding/hex"
   "os"
   "testing"
   "bytes"
)

func Test_Decrypt(t *testing.T) {
   enc, err := os.ReadFile("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   dec, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dec.Close()
   key, err := hex.DecodeString("680a46ebd6cf2b9a6a0b05a24dcf944a")
   if err != nil {
      t.Fatal(err)
   }
   sinf, err := decrypt(bytes.NewReader(enc), dec)
   if err != nil {
      t.Fatal(err)
   }
   if err := sinf.decrypt(bytes.NewReader(enc), key, dec); err != nil {
      t.Fatal(err)
   }
}
