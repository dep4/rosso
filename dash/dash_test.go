package dash

import (
   "encoding/xml"
   "fmt"
   "os"
   "testing"
)

func Test_Ext(t *testing.T) {
   for _, name := range tests {
      file, err := os.Open(name)
      if err != nil {
         t.Fatal(err)
      }
      var med Media
      if err := xml.NewDecoder(file).Decode(&med); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
      fmt.Println(name)
      for _, rep := range med.Representations() {
         fmt.Printf("%q\n", rep.Ext())
      }
      fmt.Println()
   }
}
