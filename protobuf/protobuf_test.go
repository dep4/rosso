package protobuf

import (
   "encoding/json"
   "fmt"
   "os"
   "testing"
)

func TestCheckin(t *testing.T) {
   bProto, err := os.ReadFile("checkin.txt")
   if err != nil {
      t.Fatal(err)
   }
   checkinResponse, err := Unmarshal(bProto)
   if err != nil {
      t.Fatal(err)
   }
   bJSON, err := json.Marshal(checkinResponse)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Println(string(bJSON))
   androidID := checkinResponse.GetFixed64(7, "androidId")
   fmt.Println("androidId:", androidID)
}

func TestDetails(t *testing.T) {
   bProto, err := os.ReadFile("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   responseWrapper, err := Unmarshal(bProto)
   if err != nil {
      t.Fatal(err)
   }
   docV2 := responseWrapper.Get(1, "payload").
      Get(2, "detailsResponse").
      Get(4, "docV2")
   title := docV2.GetString(5, "title")
   fmt.Printf("title: %q\n", title)
   currency := docV2.Get(8, "offer").GetString(2, "currencyCode")
   fmt.Printf("currency: %q\n", currency)
   versionCode := docV2.Get(13, "details").
      Get(1, "appDetails").
      GetVarint(3, "versionCode")
   fmt.Println("versionCode:", versionCode)
   version := docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(4, "versionString")
   fmt.Printf("version: %q\n", version)
   date := docV2.Get(13, "details").
      Get(1, "appDetails").
      GetString(16, "uploadDate")
   fmt.Printf("date: %q\n", date)
   size := docV2.Get(13, "details").
      Get(1, "appDetails").
      Get(34, "installDetails").
      GetVarint(2, "size")
   fmt.Println("size:", size)
   download := docV2.Get(13, "details").
      Get(1, "appDetails").
      GetVarint(70, "numDownloads")
   fmt.Println("download:", download)
}
