package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestProtobuf(t *testing.T) {
   names := []string{"checkin.txt", "details.txt"}
   for _, name := range names {
      buf, err := os.ReadFile(name)
      if err != nil {
         t.Fatal(err)
      }
      mes, err := Unmarshal(buf)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(mes)
   }
}
