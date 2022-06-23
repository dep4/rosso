package protobuf

import (
   "bufio"
   "bytes"
   "github.com/89z/format/protobuf"
   "os"
   "testing"
)

func Benchmark_Decode(b *testing.B) {
   buf, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      b.Fatal(err)
   }
   for n := 0; n < b.N; n++ {
      _, err := Decode(bufio.NewReader(bytes.NewReader(buf)))
      if err != nil {
         b.Fatal(err)
      }
   }
}

func Benchmark_Unmarshal(b *testing.B) {
   buf, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      b.Fatal(err)
   }
   for n := 0; n < b.N; n++ {
      _, err := protobuf.Unmarshal(buf)
      if err != nil {
         b.Fatal(err)
      }
   }
}
