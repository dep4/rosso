package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func Test_Add(t *testing.T) {
   mes := make(Message)
   fmt.Println(mes.Add_Varint(1, 2))
   fmt.Println(mes.Add_String(1, "three"))
   fmt.Println(mes.Add_Varint(1, 4))
   fmt.Println(mes.Add_String(1, "five"))
   fmt.Println(mes.Add_Varint(1, 6))
   fmt.Println(mes)
}

func Test_Unmarshal(t *testing.T) {
   buf, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   response_wrapper, err := Unmarshal(buf)
   if err != nil {
      t.Fatal(err)
   }
   doc_V2 := response_wrapper.Get(1).Get(2).Get(4)
   if v := doc_V2.Get(13).Get(1).Get_Messages(17); len(v) != 4 {
      t.Fatal("File", v)
   }
   if v, err := doc_V2.Get(13).Get(1).Get_Varint(3); err != nil {
      t.Fatal(err)
   } else if v != 10218030 {
      t.Fatal("VersionCode", v)
   }
   if v, err := doc_V2.Get(13).Get(1).Get_String(4); err != nil {
      t.Fatal(err)
   } else if v != "10.21.0" {
      t.Fatal("VersionString", v)
   }
   if v, err := doc_V2.Get(13).Get(1).Get_Varint(9); err != nil {
      t.Fatal(err)
   } else if v != 47705639 {
      t.Fatal("Size", v)
   }
   if v, err := doc_V2.Get(13).Get(1).Get_String(16); err != nil {
      t.Fatal(err)
   } else if v != "Jun 14, 2022" {
      t.Fatal("Date", v)
   }
   if v, err := doc_V2.Get_String(5); err != nil {
      t.Fatal(err)
   } else if v != "Pinterest" {
      t.Fatal("title", v)
   }
   if v, err := doc_V2.Get_String(6); err != nil {
      t.Fatal(err)
   } else if v != "Pinterest" {
      t.Fatal("creator", v)
   }
   if v, err := doc_V2.Get(8).Get_String(2); err != nil {
      t.Fatal(err)
   } else if v != "USD" {
      t.Fatal("currencyCode", v)
   }
   if v, err := doc_V2.Get(13).Get(1).Get_Varint(70); err != nil {
      t.Fatal(err)
   } else if v != 750510010 {
      t.Fatal("NumDownloads", v)
   }
}
