package protobuf

import (
   "fmt"
   "os"
   "testing"
)

var buf = []byte(`
{
   "1": {"Type": 1, "Value": 10},
   "3": {
      "Type": 2, "Value": {
         "3": {"Type": 1, "Value": 11}
      }
   }
}
`)

func TestUnmarshalJSON(t *testing.T) {
   mes := make(message)
   err := mes.UnmarshalJSON(buf)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", mes)
}

func TestMarshalJSON(t *testing.T) {
   mes := unmarshal(youtube)
   buf, err := mes.MarshalJSON()
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(buf)
}
