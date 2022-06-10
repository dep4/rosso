package decrypt

import (
   "encoding/hex"
   "fmt"
   "os"
   "testing"
)

const rawKey = "6b1f79ba70956a37fe716997b8d211ae"

var segments = []string{
   "segment0.m4f",
   "segment1.m4f",
   "segment2.m4f",
   "segment3.m4f",
   "segment4.m4f",
   "segment5.m4f",
   "segment6.m4f",
   "segment7.m4f",
   "segment8.m4f",
   "segment9.m4f",
}

func TestDecrypt(t *testing.T) {
   dec, err := os.Create("ignore/dec.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dec.Close()
   init0, err := os.Open("ignore/init0.m4f")
   if err != nil {
      t.Fatal(err)
   }
   defer init0.Close()
   if err := decryptInit(init0, dec); err != nil {
      t.Fatal(err)
   }
   key, err := hex.DecodeString(rawKey)
   if err != nil {
      t.Fatal(err)
   }
   for _, segment := range segments {
      fmt.Println(segment)
      file, err := os.Open("ignore/" + segment)
      if err != nil {
         t.Fatal(err)
      }
      if err := decryptSegment(file, key, dec); err != nil {
         t.Fatal(err)
      }
      if err := file.Close(); err != nil {
         t.Fatal(err)
      }
   }
}
