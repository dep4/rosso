package protobuf

import (
   "os"
   "testing"
)

func TestCheckin(t *testing.T) {
   file, err := os.Open("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   responseWrapper, err := Decode(file)
   if err != nil {
      t.Fatal(err)
   }
   docV2 := responseWrapper.Get(1).Get(2).Get(4)
   if v, err := docV2.Get(13).Get(1).GetVarint(3); err != nil {
      t.Fatal(err)
   } else if v != 10218030 {
      t.Fatal("VersionCode", v)
   }
   if v, err := docV2.Get(13).Get(1).GetString(4); err != nil {
      t.Fatal(err)
   } else if v != "10.21.0" {
      t.Fatal("VersionString", v)
   }
   if v, err := docV2.Get(13).Get(1).GetVarint(9); err != nil {
      t.Fatal(err)
   } else if v != 47705639 {
      t.Fatal("Size", v)
   }
   if v, err := docV2.Get(13).Get(1).GetString(16); err != nil {
      t.Fatal(err)
   } else if v != "Jun 14, 2022" {
      t.Fatal("Date", v)
   }
   if v := docV2.Get(13).Get(1).GetMessages(17); len(v) != 4 {
      t.Fatal("file", v)
   }
   if v, err := docV2.GetString(5); err != nil {
      t.Fatal(err)
   } else if v != "Pinterest" {
      t.Fatal("title", v)
   }
   if v, err := docV2.GetString(6); err != nil {
      t.Fatal(err)
   } else if v != "Pinterest" {
      t.Fatal("creator", v)
   }
   if v, err := docV2.Get(8).GetString(2); err != nil {
      t.Fatal(err)
   } else if v != "USD" {
      t.Fatal("currencyCode", v)
   }
   if v, err := docV2.Get(13).Get(1).GetVarint(70); err != nil {
      t.Fatal(err)
   } else if v != 750510010 {
      t.Fatal("NumDownloads", v)
   }
}
