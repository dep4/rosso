package protobuf

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestMarshal(t *testing.T) {
   buf := []byte("Instagram")
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%v %q\n", mes, mes.Marshal())
}

func TestUnmarshal(t *testing.T) {
   buf, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(mes)
}
