package protobuf

import (
   "bytes"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "testing"
)

func TestNew(t *testing.T) {
   b := protowire.AppendFixed64(nil, 999)
   v, err := consumeFixed64(bytes.NewReader(b))
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(v)
}
