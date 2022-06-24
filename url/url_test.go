package url

import (
   "os"
   "testing"
)

var tests = []string{
   "ignore/eol.txt",
   "ignore/noeol.txt",
}

func Test_Values(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      val, err := Decode(file)
      if err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      if err := Encode(os.Stdout, val); err != nil {
         t.Fatal(err)
      }
   }
}
