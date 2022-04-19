package xml

import (
   "fmt"
   "os"
   "testing"
)

type Form struct {
   Input []struct {
      Name string `xml:"name,attr"`
      Value string `xml:"value,attr"`
   } `xml:"input"`
}

func TestXML(t *testing.T) {
   file, err := os.Open("login.php")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   var (
      form Form
      sep = []byte(`<div class="t">`)
   )
   if err := Decode(file, sep, &form); err != nil {
      t.Fatal(err)
   }
   for _, input := range form.Input {
      fmt.Printf("%+v\n", input)
   }
}
