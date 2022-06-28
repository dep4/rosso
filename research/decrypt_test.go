package dash

import (
   "encoding/hex"
   "os"
   "testing"
)

type test_type struct {
   key string
   enc string
   dec string
}

var tests = []test_type{
   {
      "22bdb0063805260307ee5045c0f3835a",
      "ignore/enc-cbcs.mp4", "ignore/dec-cbcs.mp4",
   },
   {
      "680a46ebd6cf2b9a6a0b05a24dcf944a",
      "ignore/enc-piff.mp4", "ignore/dec-piff.mp4",
   },
}

func Test_Decrypt(t *testing.T) {
   for _, test := range tests {
      enc, err := os.Open(test.enc)
      if err != nil {
         t.Fatal(err)
      }
      defer enc.Close()
      dec, err := os.Create(test.dec)
      if err != nil {
         t.Fatal(err)
      }
      defer dec.Close()
      key, err := hex.DecodeString(test.key)
      if err != nil {
         t.Fatal(err)
      }
      if err := decryptMP4withCenc(enc, key, dec); err != nil {
         t.Fatal(err)
      }
   }
}
