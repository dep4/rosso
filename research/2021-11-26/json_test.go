package protobuf

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

var buf = []byte(`
{
   "1": {"Type": 1, "Value": 10},
   "2": {
      "Type": 2, "Value": {
         "3": {"Type": 1, "Value": 11}
      }
   },
   "3": [
      {"Type": 1, "Value": 10}, {"Type": 1, "Value": 11}
   ]
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


type appDetails struct {
   Version []varint `json:"1"`
}

func TestVarint(t *testing.T) {
   app := appDetails{
      []varint{2, 3},
   }
   // {"1":[{"Type":0,"Value":2},{"Type":0,"Value":3}]}
   err := json.NewEncoder(os.Stdout).Encode(app)
   if err != nil {
      t.Fatal(err)
   }
}
