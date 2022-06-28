package dash

import (
   "bytes"
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

var tests = []test_type{
   // pass
   {
      "680a46ebd6cf2b9a6a0b05a24dcf944a",
      "ignore/enc-piff.mp4", "ignore/dec-piff.mp4",
   },
   // fail
   {
      "22bdb0063805260307ee5045c0f3835a",
      "ignore/enc-cbcs.mp4", "ignore/dec-cbcs.mp4",
   },
}

func Test_Two(t *testing.T) {
   for _, test := range tests {
      fmt.Println(test.enc)
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
