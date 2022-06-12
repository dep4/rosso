package json

import (
   "fmt"
   "os"
   "testing"
)

func TestScanner(t *testing.T) {
   file, err := os.Open("roku.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   var scan Scanner
   if _, err := scan.ReadFrom(file); err != nil {
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
