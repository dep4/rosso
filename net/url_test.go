package net

import (
   "os"
   "testing"
)

var tests = []string{
   "ignore/eol.txt",
   "ignore/noeol.txt",
}

func TestValues(t *testing.T) {
   for _, test := range tests {
      file, err := os.Open(test)
      if err != nil {
         t.Fatal(err)
      }
      val := NewValues()
      if _, err := val.ReadFrom(file); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      if _, err := val.WriteTo(os.Stdout); err != nil {
         t.Fatal(err)
      }
   }
}
