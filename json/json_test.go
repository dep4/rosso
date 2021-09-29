package json

import (
   "fmt"
   "os"
   "testing"
)

func TestJSON(t *testing.T) {
   data, err := os.ReadFile("ig.js")
   if err != nil {
      t.Fatal(err)
   }
   var v struct {
      Shortcode_Media struct {
         ID string
      }
   }
   if err := Unmarshal(data, &v); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", v)
}
