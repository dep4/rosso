package dash

import (
   "bytes"
   "encoding/hex"
   "os"
   "testing"
)

const (
   rawKey = "22bdb0063805260307ee5045c0f3835a"
   zero = "P377684155_A1524726231_FF_video_gr203_sdr_508x254_cbcs_--0.mp4"
   one = "P377684155_A1524726231_FF_video_gr203_sdr_508x254_cbcs_--1.m4s"
)

func TestDecrypt(t *testing.T) {
   dec, err := os.Create("ignore/dec.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer dec.Close()
   zero, err := os.ReadFile("ignore/" + zero)
   if err != nil {
      t.Fatal(err)
   }
   dec.Write(zero)
   one, err := os.ReadFile("ignore/" + one)
   if err != nil {
      t.Fatal(err)
   }
   buf := new(bytes.Buffer)
   buf.Write(zero)
   buf.Write(one)
   key, err := hex.DecodeString(rawKey)
   if err != nil {
      t.Fatal(err)
   }
   if err := Decrypt(dec, buf, key); err != nil {
      t.Fatal(err)
   }
}
