package protobuf

import (
   "fmt"
   "os"
   "testing"
)

// string first:
var files = []string{
   // {4 }:16.49.37 
   "com.google.android.youtube.txt",
   // {4 }:216.1.0.21.137
   "com.instagram.android.txt",
}

func TestProto(t *testing.T) {
   for _, file := range files {
      buf, err := os.ReadFile(file)
      if err != nil {
         t.Fatal(err)
      }
      mes, err := Unmarshal(buf)
      if err != nil {
         t.Fatal(err)
      }
      appDetails := mes.Get(1, 2, 4, 13, 1)
      fmt.Print(file, ":\n", appDetails, "\n")
   }
}
