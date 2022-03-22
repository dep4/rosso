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
   fmt.Println(mes.GetUint64(7))
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
   fmt.Println(len(buf), len(mes.Marshal()))
   fmt.Printf("%q\n", mes.Get(1).Get(2).Get(4).GetString(5))
   for _, image := range mes.Get(1).Get(2).Get(4).GetMessages(10) {
      fmt.Println(image)
   }
}
