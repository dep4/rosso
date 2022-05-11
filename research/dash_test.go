package dash

import (
   "encoding/hex"
   "os"
   "testing"
)

const hexKey = "13d7c7cf295444944b627ef0ad2c1b3c"

func TestDecrypt(t *testing.T) {
   dst, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dst.Close()
   init, err := os.Open("ignore/init.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer init.Close()
   dst.ReadFrom(init)
   media, err := os.Open("ignore/enc.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer media.Close()
   key, err := hex.DecodeString(hexKey)
   if err != nil {
      t.Fatal(err)
   }
   if err := decryptMP4withCenc(media, key, dst); err != nil {
      t.Fatal(err)
   }
}
