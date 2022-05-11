package dash

import (
   "encoding/hex"
   "io"
   "os"
   "testing"
)

const hexKey = "13d7c7cf295444944b627ef0ad2c1b3c"

func write(dst io.Writer, key []byte, name string) error {
   file, err := os.Open(name)
   if err != nil {
      return err
   }
   defer file.Close()
   return decryptMP4withCenc(file, key, dst)
}

func TestRoku(t *testing.T) {
   dst, err := os.Create("ignore.mp4")
   if err != nil {
      t.Fatal(err)
   }
   init, err := os.Open("ignore/init.mp4")
   if err != nil {
      t.Fatal(err)
   }
   defer init.Close()
   dst.ReadFrom(init)
   key, err := hex.DecodeString(hexKey)
   if err != nil {
      t.Fatal(err)
   }
   if err := write(dst, key, "ignore/1.mp4"); err != nil {
      t.Fatal(err)
   }
   if err := write(dst, key, "ignore/2.mp4"); err != nil {
      t.Fatal(err)
   }
}
