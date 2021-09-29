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
   // array
   var a []struct {
      Src string
   }
   if err := UnmarshalArray(data, &a); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", a)
   // object
   var o struct {
      Shortcode_Media struct {
         ID string
      }
   }
   if err := UnmarshalObject(data, &o); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", o)
}
