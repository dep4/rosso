package protobuf

import (
   "encoding/json"
   "os"
   "testing"
)

func TestProto(t *testing.T) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   toks := ParseUnknown(data)
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   if err := enc.Encode(toks); err != nil {
      t.Fatal(err)
   }
}
