package xml

import (
   "fmt"
   "os"
   "testing"
)

func TestInput(t *testing.T) {
   file, err := os.Open("login.php")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte("<input ")
   for scan.Scan() {
      var input struct {
         Name string `xml:"name,attr"`
         Value string `xml:"value,attr"`
      }
      err := scan.Decode(&input)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", input)
   }
}

func TestMeta(t *testing.T) {
   file, err := os.Open("ignore.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte("<meta ")
   for scan.Scan() {
      var meta struct {
         Property string `xml:"property,attr"`
         Content string `xml:"content,attr"`
      }
      err := scan.Decode(&meta)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%+v\n", meta)
   }
}
