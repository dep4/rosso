package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestTwo(t *testing.T) {
   buf := []byte("Instagram")
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%v %q\n", mes, mes.Marshal())
}

func TestOne(t *testing.T) {
   buf, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(mes)
}
