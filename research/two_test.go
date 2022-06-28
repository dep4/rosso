package dash

import (
   "encoding/hex"
   "os"
   "testing"
   "bytes"
)

func Test_Two(t *testing.T) {
   for _, test := range tests {
      file, err := os.Create(test.dec)
      if err != nil {
         t.Fatal(err)
      }
      defer file.Close()
      dec := new_decrypter(file)
      buf, err := os.ReadFile(test.enc)
      if err != nil {
         t.Fatal(err)
      }
      if err := dec.init(bytes.NewReader(buf)); err != nil {
         t.Fatal(err)
      }
      key, err := hex.DecodeString(test.key)
      if err != nil {
         t.Fatal(err)
      }
      if err := dec.segment(bytes.NewReader(buf), key); err != nil {
         t.Fatal(err)
      }
   }
}
