package protobuf

import (
   "encoding/json"
   "os"
   "testing"
)

func TestCheckin(t *testing.T) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes := make(Message)
   if err := mes.UnmarshalBinary(data); err != nil {
      t.Fatal(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(mes)
}
