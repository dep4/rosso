package dash

import (
   "encoding/hex"
   "os"
   "testing"
   "bytes"
)

func Test_Decrypt(t *testing.T) {
   out, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer out.Close()
   dec := new_decrypter(out)
   in, err := os.ReadFile("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   if err := dec.init(bytes.NewReader(in)); err != nil {
      t.Fatal(err)
   }
   key, err := hex.DecodeString("680a46ebd6cf2b9a6a0b05a24dcf944a")
   if err != nil {
      t.Fatal(err)
   }
   if err := dec.segment(bytes.NewReader(in), key); err != nil {
      t.Fatal(err)
   }
}
