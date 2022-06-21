package protobuf

import (
   "encoding/json"
   "os"
   "testing"
)

func Test_Checkin(t *testing.T) {
   buf, err := os.ReadFile("com.pinterest.txt")
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
