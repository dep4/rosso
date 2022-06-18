package protobuf

import (
   "encoding/json"
   "os"
   "testing"
)

func TestUnmarshal(t *testing.T) {
   buf, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes := make(Message)
   if err := mes.UnmarshalBinary(buf); err != nil {
      t.Fatal(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(mes)
}
