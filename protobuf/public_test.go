package protobuf

import (
   "fmt"
   "os"
   "testing"
)

var tests = []string{
   "checkin.txt",
   "details.txt",
}

func TestCheckin(t *testing.T) {
   for _, test := range tests {
      buf, err := os.ReadFile(test)
      if err != nil {
         t.Fatal(err)
      }
      mes, err := Unmarshal(buf)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Println(len(mes.Marshal()))
   }
}
