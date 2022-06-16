package protobuf

import (
   "bufio"
   "fmt"
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
   creator, err := docV2.GetString(6)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", creator)
   currencyCode, err := docV2.Get(8).GetString(2)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%q\n", currencyCode)
   /*
   det.Micros, err = docV2.Get(8).GetVarint(1)
   det.NumDownloads, err = docV2.Get(13).Get(1).GetVarint(70)
   det.Size, err = docV2.Get(13).Get(1).GetVarint(9)
   det.Title, err = docV2.GetString(5)
   det.UploadDate, err = docV2.Get(13).Get(1).GetString(16)
   det.VersionCode, err = docV2.Get(13).Get(1).GetVarint(3)
   det.VersionString, err = docV2.Get(13).Get(1).GetString(4)
   for _, file := range docV2.Get(13).Get(1).GetMessages(17) {
      typ, err := file.GetVarint(1)
   }
   */
}
