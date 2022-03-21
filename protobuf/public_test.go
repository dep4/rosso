package protobuf

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestCheckin(t *testing.T) {
   buf, err := os.ReadFile("checkin.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(mes)
   id := Value[uint64](mes, 7)
   fmt.Println(id)
}

func TestDetails(t *testing.T) {
   buf, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   mes, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   enc := json.NewEncoder(os.Stdout)
   enc.SetIndent("", " ")
   enc.Encode(mes)
   title := Value[string](mes, 1, 2, 4, 5)
   fmt.Printf("%q\n", title)
   for _, image := range Values[Message](mes, 1, 2, 4, 10) {
      fmt.Println(image)
   }
}
