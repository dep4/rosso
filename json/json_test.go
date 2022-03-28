package json

import (
   "fmt"
   "os"
   "testing"
)

type sigiData struct {
   ItemModule map[int]struct {
      Author string
      ID string
      Video struct {
         PlayAddr string
      }
   }
}

func TestJSON(t *testing.T) {
   buf, err := os.ReadFile("tiktok.js")
   if err != nil {
      t.Fatal(err)
   }
   var data sigiData
   ok := NewDecoder(buf).Object(&data)
   fmt.Printf("%+v %v\n", data, ok)
}
