package json

import (
   "fmt"
   "os"
   "testing"
)

func TestJSON(t *testing.T) {
   buf, err := os.ReadFile("eD41U.json")
   if err != nil {
      t.Fatal(err)
   }
   var v []byte
   dec := NewDecoder(buf)
   ok := dec.DecodeArray(&v)
   fmt.Println(ok, v)
   ok = dec.DecodeArray(&v)
   fmt.Println(ok, v)
}
