package json

import (
   "fmt"
   "os"
   "testing"
)

/*
mech\pbs\frontline.go
mech\pbs\masterpiece.go
mech\pbs\nature.go
mech\pbs\nova.go
mech\pbs\video.go
*/
func TestPBS(t *testing.T) {
   file, err := os.Open("widget.html")
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
   var video struct {
      Encodings []string
   }
   if err := scan.Decode(&video); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", video)
}

func TestFacebook(t *testing.T) {
   file, err := os.Open("ignore.html")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   scan, err := NewScanner(file)
   if err != nil {
      t.Fatal(err)
   }
   scan.Split = []byte(`{"`)
   scan.Scan()
   var post struct {
      DateCreated string
   }
   if err := scan.Decode(&post); err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%+v\n", post)
}
