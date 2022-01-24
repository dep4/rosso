package protobuf

import (
   "fmt"
   "os"
   "testing"
)

var apps = []string{
   "com.google.android.youtube.txt",
   "com.instagram.android.txt",
}

func TestBinary(t *testing.T) {
   for _, app := range apps {
      buf, err := os.ReadFile(app)
      if err != nil {
         t.Fatal(err)
      }
      responseWrapper, err := Unmarshal(buf)
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
         GetUint64(3, "versionCode")
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
         GetUint64(2, "size")
      fmt.Println("size:", size)
      download := docV2.Get(13, "details").
         Get(1, "appDetails").
         GetUint64(70, "numDownloads")
      fmt.Println("download:", download)
   }
}
