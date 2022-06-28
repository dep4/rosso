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

func Test_One(t *testing.T) {
   for _, test := range tests {
      file, err := os.Create(test.dec)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      key, err := hex.DecodeString(test.key)
      if err != nil {
         t.Fatal(err)
      }
      enc, err := os.Open(test.enc)
      if err != nil {
         t.Fatal(err)
      }
      defer enc.Close()
      if err := decryptMP4withCenc(enc, key, file); err != nil {
         t.Fatal(err)
      }
   }
}
