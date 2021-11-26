package protobuf

import (
   "testing"
   "os"
)

func TestMarshalJSON(t *testing.T) {
   mes := unmarshal(youtube)
   buf, err := mes.MarshalJSON()
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(buf)
}
