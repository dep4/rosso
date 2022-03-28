package json

import (
   "fmt"
   "os"
   "testing"
)

func TestJSON(t *testing.T) {
   buf, err := os.ReadFile("3016754074.html")
   if err != nil {
      t.Fatal(err)
   }
   var videoBridge struct {
      Encodings []string
   }
   var bufs = [][]byte{nil, buf}
   var seps = [][]byte{
      nil, []byte(" = "), []byte("\twindow.videoBridge = "),
   }
   for _, buf := range bufs{
      for _, sep := range seps {
         err := Unmarshal(buf, sep, &videoBridge)
         fmt.Printf("%+v %v\n", videoBridge, err)
      }
   }
}
