package protobuf

import (
   "bufio"
   "os"
   "testing"
)

func TestCheckin(t *testing.T) {
   file, err := os.Open("details.txt")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   responseWrapper, err := readMessage(bufio.NewReader(file))
   if err != nil {
      t.Fatal(err)
   }
   docV2 := responseWrapper.Get(1).Get(2).Get(4)
   if v, err := docV2.GetString(6); err != nil {
      t.Fatal(err)
   } else if v != "Instagram" {
      t.Fatal(v)
   }
   if v, err := docV2.Get(8).GetString(2); err != nil {
      t.Fatal(err)
   } else if v != "USD" {
      t.Fatal(v)
   }
   if v, err := docV2.Get(13).Get(1).GetVarint(70); err != nil {
      t.Fatal(err)
   } else if v != 3931864786 {
      t.Fatal(err)
   }
   if v, err := docV2.Get(13).Get(1).GetVarint(9); err != nil {
      t.Fatal(err)
   } else if v != 52627455 {
      t.Fatal(v)
   }
   if v, err := docV2.GetString(5); err != nil {
      t.Fatal(err)
   } else if v != "Instagram" {
      t.Fatal(v)
   }
   if v, err := docV2.Get(13).Get(1).GetString(16); err != nil {
      t.Fatal(err)
   } else if v != "Dec 16, 2021" {
      t.Fatal(v)
   }
   if v, err := docV2.Get(13).Get(1).GetVarint(3); err != nil {
      t.Fatal(err)
   } else if v != 321704040 {
      t.Fatal(v)
   }
   if v, err := docV2.Get(13).Get(1).GetString(4); err != nil {
      t.Fatal(err)
   } else if v != "216.1.0.21.137" {
      t.Fatal(v)
   }
   files := docV2.Get(13).Get(1).GetMessages(17)
   if len(files) != 1 {
      t.Fatal(files)
   }
}
