package json

import (
   "bytes"
   "os"
   "testing"
)

func BenchmarkNew(b *testing.B) {
   buf, err := os.ReadFile("ignore.html")
   if err != nil {
      b.Fatal(err)
   }
   for n := 0; n < b.N; n++ {
      NewScanner(bytes.NewReader(buf))
   }
}

func BenchmarkRead(b *testing.B) {
   buf, err := os.ReadFile("ignore.html")
   if err != nil {
      b.Fatal(err)
   }
   var scan Scanner
   for n := 0; n < b.N; n++ {
      scan.ReadFrom(bytes.NewReader(buf))
   }
}
