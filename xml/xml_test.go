package xml

import (
   "fmt"
   "os"
   "testing"
)

func TestXML(t *testing.T) {
   file, err := os.Open("ignore.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   var meta struct {
      Content string `xml:"content,attr"`
   }
   scan.Split = []byte(`"web-tv-app/config/environment"`)
   scan.Scan()
   scan.Split = []byte("<meta")
   if err := scan.Decode(&meta); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", meta)
}
