package net

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{"meta", "script"}

func TestMeta(t *testing.T) {
   for _, test := range tests {
      f, err := os.Open("bleep.html")
      if err != nil {
         t.Fatal(err)
      }
      defer f.Close()
      for _, node := range ReadHTML(f, test) {
         fmt.Printf("%+v\n", node)
      }
   }
}
