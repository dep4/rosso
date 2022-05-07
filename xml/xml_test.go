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
   var script struct {
      DataTralbum []byte `xml:"data-tralbum,attr"`
   }
   scan.Split = []byte(" data-tralbum=")
   scan.Scan()
   scan.Split = []byte("<script data-tralbum=")
   if err := scan.Decode(&script); err != nil {
      t.Fatal(err)
   }
   fmt.Println(string(script.DataTralbum))
}
