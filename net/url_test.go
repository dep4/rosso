package net

import (
   "os"
   "testing"
)

type value_test struct {
   in string
   out int64
}

var tests = []value_test{
   {"ignore/eol.txt", 446},
   {"ignore/noeol.txt", 679},
}

func Test_Values(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test.in)
      if err != nil {
         t.Fatal(err)
      }
      val := New_Values()
      num, err := val.ReadFrom(file)
      if err != nil {
         t.Fatal(err)
      }
      if num != test.out {
         t.Fatal(num)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      if _, err := val.WriteTo(os.Stdout); err != nil {
         t.Fatal(err)
      }
   }
}
