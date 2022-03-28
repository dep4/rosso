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
   sep := []byte("\twindow.videoBridge = ")
   if err := Unmarshal(buf, sep, &videoBridge); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", videoBridge)
}
