package json

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func TestDecode(t *testing.T) {
   src := strings.NewReader(`{"month":12,"day":31}`)
   var date struct {
      Month int
      Day int
   }
   err := NewDecoder(src).Decode(&date)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", date)
}

func TestPBS(t *testing.T) {
   file, err := os.Open("pbs-widget.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"availability"`)
   scan.Scan()
   var object struct {
      Encodings []string
   }
   if err := scan.Decode(&object); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", object)
}
