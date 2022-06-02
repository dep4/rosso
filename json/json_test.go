package json

import (
   "fmt"
   "os"
   "strings"
   "testing"
)

func TestScanner(t *testing.T) {
   file, err := os.Open("roku.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte("\tcsrf:")
   scan.Scan()
   scan.Split = nil
   var token string
   if err := scan.Decode(&token); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", token)
}

func TestDecoder(t *testing.T) {
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
