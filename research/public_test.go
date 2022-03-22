package protobuf

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestDetails(t *testing.T) {
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
   fmt.Println(len(buf), len(mes.Marshal()))
}
