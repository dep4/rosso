package protobuf

import (
   "encoding/json"
   "os"
   "testing"
)

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
