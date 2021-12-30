package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestProto(t *testing.T) {
   // string first:
   // {4 }:9.0.12.50
   // message first:
   // {4 }:map[{7 }:3473733480694558766]
   buf, err := os.ReadFile("com.sec.android.app.launcher.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(mes)
}
