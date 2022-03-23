package protobuf

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestJSON(t *testing.T) {
   buf := []byte(`{"month": 12, "day": 31}`)
   { // pass
      var date struct { Month int }
      err := json.Unmarshal(buf, &date)
      fmt.Printf("%+v %v\n", date, err)
   }
   { // fail key
      var date struct { Year int }
      err := json.Unmarshal(buf, &date)
      fmt.Printf("%+v %v\n", date, err)
   }
   { // fail type
      var date struct { Month string }
      err := json.Unmarshal(buf, &date)
      fmt.Printf("%+v %v\n", date, err)
   }
}

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
   // pass
   fmt.Println(mes.GetFixed64(7))
   // fail key
   fmt.Println(mes.GetFixed64(6))
   // fail type
   fmt.Println(mes.GetVarint(7))
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
   // .payload.detailsResponse.docV2
   docV2 := mes.Get(1).Get(2).Get(4)
   // .title
   title, ok := docV2.GetString(5)
   fmt.Printf("%q %v\n", title, ok)
   // .details.appDetails.file
   for _, file := range docV2.Get(13).Get(1).GetMessages(17) {
      fmt.Println(file)
   }
}
